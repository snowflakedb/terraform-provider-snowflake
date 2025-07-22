package testfunctional

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

type ObjectRenamingDatabaseListItem struct {
	String string
	Int    int
}

func mapObjectRenamingDatabaseListItemFromValue(items []cty.Value) []ObjectRenamingDatabaseListItem {
	return collections.Map(items, func(item cty.Value) ObjectRenamingDatabaseListItem {
		intValue, _ := item.AsValueMap()["int"].AsBigFloat().Int64()
		return ObjectRenamingDatabaseListItem{
			String: item.AsValueMap()["string"].AsString(),
			Int:    int(intValue),
		}
	})
}

func objectRenamingDatabaseListFromSchema(items []any) []ObjectRenamingDatabaseListItem {
	return collections.Map(items, func(item any) ObjectRenamingDatabaseListItem {
		return ObjectRenamingDatabaseListItem{
			String: item.(map[string]any)["string"].(string),
			Int:    item.(map[string]any)["int"].(int),
		}
	})
}

type objectRenamingDatabaseOrderedListItem struct {
	Name  string
	Order string
}

func mapObjectRenamingDatabaseOrderedListItemFromValue(items []cty.Value) []objectRenamingDatabaseOrderedListItem {
	return collections.Map(items, func(item cty.Value) objectRenamingDatabaseOrderedListItem {
		var order string
		if orderValue, ok := item.AsValueMap()["order"]; ok && !orderValue.IsNull() {
			order = orderValue.AsString()
		}
		var name string
		if nameValue, ok := item.AsValueMap()["name"]; ok && !nameValue.IsNull() {
			name = nameValue.AsString()
		}
		return objectRenamingDatabaseOrderedListItem{
			Name:  name,
			Order: order,
		}
	})
}

func objectRenamingDatabaseOrderedListFromSchema(list []any) []objectRenamingDatabaseOrderedListItem {
	objectRenamingDatabaseOrderedListItems := make([]objectRenamingDatabaseOrderedListItem, len(list))
	for index, item := range list {
		var name string
		if nameValue, ok := item.(map[string]any)["name"]; ok {
			name = nameValue.(string)
		}
		objectRenamingDatabaseOrderedListItems[index] = objectRenamingDatabaseOrderedListItem{
			Name:  name,
			Order: strconv.Itoa(index),
		}
	}
	return objectRenamingDatabaseOrderedListItems
}

type ObjectRenamingDatabaseManuallyOrderedListItem struct {
	Name string
	Type string
}

func objectRenamingDatabaseManuallyOrderedListFromSchema(list []any) []ObjectRenamingDatabaseManuallyOrderedListItem {
	objectRenamingDatabaseOrderedListItems := make([]ObjectRenamingDatabaseManuallyOrderedListItem, len(list))
	slices.SortFunc(list, func(a, b any) int {
		return a.(map[string]any)["order"].(int) - b.(map[string]any)["order"].(int)
	})
	for index, item := range list {
		objectRenamingDatabaseOrderedListItems[index] = ObjectRenamingDatabaseManuallyOrderedListItem{
			Name: item.(map[string]any)["name"].(string),
			Type: item.(map[string]any)["type"].(string),
		}
	}
	return objectRenamingDatabaseOrderedListItems
}

type objectRenamingDatabase struct {
	List                []ObjectRenamingDatabaseListItem
	OrderedList         []objectRenamingDatabaseOrderedListItem
	ManuallyOrderedList []ObjectRenamingDatabaseManuallyOrderedListItem
	ChangeLog           ObjectRenamingDatabaseChangelog
}

type ObjectRenamingDatabaseChangelogChange struct {
	Before map[string]any
	After  map[string]any
}

// ObjectRenamingDatabaseChangelog is used for testing purposes to track actions taken in the Update method like Add/Remove/Change.
// It's only supported the manually_ordered_list option.
type ObjectRenamingDatabaseChangelog struct {
	Added   []map[string]any
	Removed []map[string]any
	Changed []ObjectRenamingDatabaseChangelogChange
}

var ObjectRenamingDatabaseInstance = &objectRenamingDatabase{
	List:                make([]ObjectRenamingDatabaseListItem, 0),
	OrderedList:         make([]objectRenamingDatabaseOrderedListItem, 0),
	ManuallyOrderedList: make([]ObjectRenamingDatabaseManuallyOrderedListItem, 0),
	ChangeLog:           ObjectRenamingDatabaseChangelog{},
}

// TODO [SNOW-2054208]: extract to the dedicated package
var objectRenamingListsAndSetsSchema = map[string]*schema.Schema{
	// The list field was tested to be used in places where the order of the items should be ignored.
	// It was ignored by comparing hashes of the items to see if any changes were made on the items themselves
	// (if the hashes before and after were the same, we know that nothing was changed, only the order).
	// Also, it doesn't fully support repeating items. This is because they have the same hash and to fully support it,
	// hash counting could be added (counting if the same hash occurs in state and config the number of times, otherwise cause update).
	// Modifications of the items will still cause remove/add behavior.
	"list": {
		Optional:         true,
		Type:             schema.TypeList,
		DiffSuppressFunc: ignoreListOrderAfterFirstApply("list"),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"string": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"int": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	},
	// The manually_ordered_list focuses on providing both aspects:
	//  - Immunity to item reordering after create.
	//  - Handling updates for changed items instead of removing the old item and adding a new one.
	// It does it by providing the required order field that represents what should be the actual order of items
	// on the Snowflake side. The order is ignored on the DiffSuppressFunc level, and the item update (renaming the item)
	// is handled in the resource update function. This proposal is supposed to test the behavior of Snowflake columns needed for table refactor for v1.
	// Here's the full list of what should be possible with this approach:
	// Supported actions:
	// - Drop item (any position).
	// - Add item (at the end).
	// - Rename item / Change item type (any position; compatible type change).
	// - Reorder items.
	// Unsupported actions:
	// - Add item (in the middle).
	// - Change item type (incompatible change).
	// - External changes (with an option to set either ForceNew or error behavior).
	// Assumptions:
	// - The list "returned from Snowflake side" is ordered (or identifiable).
	// - Order field is treated as an identifier that cannot be changed for the lifetime of a given item.
	// - Items contain fields that are able to uniquely identify a given item (in this case, we have name + type).
	"manually_ordered_list": {
		Optional:         true,
		Type:             schema.TypeList,
		DiffSuppressFunc: ignoreOrderAfterFirstApplyWithManuallyOrderedList("manually_ordered_list"),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"order": {
					Type:     schema.TypeInt,
					Required: true,
					// Improvement:
					// Cause ForceNew behavior whenever any of the items change their order to different value than previously.
					// It's not trivial as it cannot be achieved by putting ForceNew modifier in the schema (with the current implementation of Update/Read/SuppressDiff).
					// It also cannot be achieved by creating custom diff. It seems the custom diff is seeing the changes
					// too late to call ForceNew and for Terraform to show it during the plan or apply it during the apply.
					// Currently, the only good way to prevent such changes is to describe them clearly in the documentation.
				},
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateDiagFunc: resources.SdkValidation(func(value string) (string, error) {
						if slices.Contains([]string{"INT", "NUMBER", "STRING", "TEXT"}, value) {
							return value, nil
						}
						return "", fmt.Errorf("invalid type: %s", value)
					}),
				},
			},
		},
	},
	// The invalid_operation_handler is a switch that we wanted to implement for the approach with manually_ordered_list.
	// The idea behind it was to switch between different error handling behaviors for invalid operations we discussed
	// (see unsupported actions in the manually_ordered_list description). For now, we wanted to have two options:
	// error out or force re-creation (force new). Currently, the proposal uses only errors because during FORCE_NEW
	// testing it seemed like it's impossible to for resource re-creation without adding more logic into custom diff
	// that would be able to calculate it on plan time, not apply time.
	"invalid_operation_handler": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "ERROR",
		ValidateDiagFunc: resources.SdkValidation(func(value string) (string, error) {
			if slices.Contains([]string{"ERROR", "FORCE_NEW"}, value) {
				return value, nil
			}
			return "", fmt.Errorf("invalid invalid operation handler: %s", value)
		}),
	},
	// The invalid_operation field was an attempt to collect error messages gathered in the Update and Read methods
	// and use it in custom diff that would try to re-create the resource on non-empty invalid_operation field.
	// The FORCE_NEW handler applied this way seemed to have no impact on the resource behavior because the information was
	// transferred to the field too late. As mentioned, it happened during Update and Read when the `terraform apply` is
	// already running. During the run, it's not valid to apply the force new, because it has to be known before the apply, so
	// Terraform would be able to show it during the plan. Because of that, we know that the logic inside custom diff has to
	// do much more guessing on its own, so we would be able to know those invalid operations during the plan time.
	// Only then we would be able to have invalid_operation_handler and FORCE_NEW as a valid option (and this field would be
	// most likely useless and could be removed).
	"invalid_operation": {
		Type:     schema.TypeString,
		Computed: true,
	},
	// The ordered_list field was an attempt of making manual work done in manually_ordered_list automatic by making the order field computed.
	// It didn't work because in DiffSuppressFunc it's hard to get the computed value in the "after" state to compare against.
	// Possibly (but with very low probability), the solution could work by introducing a Computed + Optional list that would be managed by a custom diff function.
	// Due to increased complexity, it was left as is and more research was dedicated to manually_ordered_list.
	"ordered_list": {
		Optional:         true,
		Type:             schema.TypeList,
		DiffSuppressFunc: ignoreOrderAfterFirstApplyWithOrderedList("ordered_list"),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"order": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	},
}

func TestResourceObjectRenamingListsAndSets() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateObjectRenamingListsAndSets,
		UpdateContext: UpdateObjectRenamingListsAndSets,
		ReadContext:   ReadObjectRenamingListsAndSets(true),
		DeleteContext: DeleteObjectRenamingListsAndSets,

		Schema: objectRenamingListsAndSetsSchema,
	}
}

func CreateObjectRenamingListsAndSets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ObjectRenamingDatabaseInstance.List = objectRenamingDatabaseListFromSchema(d.Get("list").([]any))
	ObjectRenamingDatabaseInstance.OrderedList = objectRenamingDatabaseOrderedListFromSchema(d.Get("ordered_list").([]any))
	ObjectRenamingDatabaseInstance.ManuallyOrderedList = objectRenamingDatabaseManuallyOrderedListFromSchema(d.Get("manually_ordered_list").([]any))

	d.SetId("identifier")

	return ReadObjectRenamingListsAndSets(false)(ctx, d, meta)
}

func UpdateObjectRenamingListsAndSets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if d.HasChange("list") {
		// It wasn't working with d.getChange(). It was returning null elements in one of the test case's steps.
		oldList := d.GetRawState().AsValueMap()["list"].AsValueSlice()
		newList := d.GetRawConfig().AsValueMap()["list"].AsValueSlice()

		oldListMapped := mapObjectRenamingDatabaseListItemFromValue(oldList)
		newListMapped := mapObjectRenamingDatabaseListItemFromValue(newList)

		addedItems, removedItems := resources.ListDiff(oldListMapped, newListMapped)

		for _, removedItem := range removedItems {
			ObjectRenamingDatabaseInstance.List = slices.DeleteFunc(ObjectRenamingDatabaseInstance.List, func(item ObjectRenamingDatabaseListItem) bool {
				shouldRemove := item == removedItem
				if shouldRemove {
					ObjectRenamingDatabaseInstance.ChangeLog.Removed = append(ObjectRenamingDatabaseInstance.ChangeLog.Removed, map[string]any{
						"int":    item.Int,
						"string": item.String,
					})
				}
				return shouldRemove
			})
		}

		for _, addedItem := range addedItems {
			ObjectRenamingDatabaseInstance.ChangeLog.Added = append(ObjectRenamingDatabaseInstance.ChangeLog.Added, map[string]any{
				"int":    addedItem.Int,
				"string": addedItem.String,
			})
		}

		ObjectRenamingDatabaseInstance.List = append(ObjectRenamingDatabaseInstance.List, addedItems...)
	}

	if d.HasChange("ordered_list") {
		oldOrderedList := d.GetRawState().AsValueMap()["ordered_list"].AsValueSlice()
		newOrderedList := d.GetRawConfig().AsValueMap()["ordered_list"].AsValueSlice()
		oldOrderedListMapped := mapObjectRenamingDatabaseOrderedListItemFromValue(oldOrderedList)
		newOrderedListMapped := mapObjectRenamingDatabaseOrderedListItemFromValue(newOrderedList)

		itemsToAdd, itemsToRemove := resources.ListDiff(oldOrderedListMapped, newOrderedListMapped)
		for _, removedItem := range itemsToRemove {
			ObjectRenamingDatabaseInstance.OrderedList = slices.DeleteFunc(ObjectRenamingDatabaseInstance.OrderedList, func(item objectRenamingDatabaseOrderedListItem) bool {
				return item == removedItem
			})
		}

		for _, addedItem := range itemsToAdd {
			addedItem.Order = ""
			ObjectRenamingDatabaseInstance.OrderedList = append(ObjectRenamingDatabaseInstance.OrderedList, addedItem)
		}

		// The implementation is not complete due to mentioned issues with computed order in DiffSuppressFunc
	}

	if d.HasChange("manually_ordered_list") {
		invalidOperations := make([]error, 0)
		updateChangelog := ObjectRenamingDatabaseChangelog{}

		oldManuallyOrderedList := d.GetRawState().AsValueMap()["manually_ordered_list"].AsValueSlice()
		newManuallyOrderedList := d.GetRawConfig().AsValueMap()["manually_ordered_list"].AsValueSlice()

		oldOrders := collections.Map(oldManuallyOrderedList, func(item cty.Value) int {
			result, _ := item.AsValueMap()["order"].AsBigFloat().Int64()
			return int(result)
		})
		maxStateOrder := slices.MaxFunc(oldOrders, func(a, b int) int { return a - b })
		finalState := make([]ObjectRenamingDatabaseManuallyOrderedListItem, 0)

		for _, oldItem := range oldManuallyOrderedList {
			oldItem := oldItem.AsValueMap()
			newItemIndex := slices.IndexFunc(newManuallyOrderedList, func(newItem cty.Value) bool {
				return oldItem["order"].AsBigFloat().Cmp(newItem.AsValueMap()["order"].AsBigFloat()) == 0
			})
			// Here we analyze already existing items and check if they need to be updated in any way.
			if newItemIndex != -1 {
				newItem := newManuallyOrderedList[newItemIndex]
				newName := oldItem["name"].AsString()
				newType := oldItem["type"].AsString()
				wasChanged := false

				if oldItem["name"].AsString() != newItem.AsValueMap()["name"].AsString() {
					// Change name
					newName = newItem.AsValueMap()["name"].AsString()
					wasChanged = true
				}

				if oldItem["type"].AsString() != newItem.AsValueMap()["type"].AsString() {
					// Change type
					newType = newItem.AsValueMap()["type"].AsString()
					wasChanged = true

					// Check for incompatible types
					if slices.Contains([]string{"TEXT", "STRING"}, oldItem["type"].AsString()) && slices.Contains([]string{"INT", "NUMBER"}, newType) ||
						slices.Contains([]string{"INT", "NUMBER"}, oldItem["type"].AsString()) && slices.Contains([]string{"TEXT", "STRING"}, newType) {
						invalidOperations = append(invalidOperations, fmt.Errorf("unable to change item type from %s to %s", oldItem["type"].AsString(), newType))
					}
				}

				itemToAdd := ObjectRenamingDatabaseManuallyOrderedListItem{
					Name: newName,
					Type: newType,
				}
				finalState = append(finalState, itemToAdd)

				if wasChanged {
					updateChangelog.Changed = append(updateChangelog.Changed, ObjectRenamingDatabaseChangelogChange{
						Before: map[string]any{
							"name": oldItem["name"].AsString(),
							"type": oldItem["type"].AsString(),
						},
						After: map[string]any{
							"name": itemToAdd.Name,
							"type": itemToAdd.Type,
						},
					})
				}
			} else {
				// If given order wasn't found, it means this item was removed.
				updateChangelog.Removed = append(updateChangelog.Removed, map[string]any{
					"name": oldItem["name"].AsString(),
					"type": oldItem["type"].AsString(),
				})
			}
		}

		// Here we analyze newly added items
		for _, newItem := range newManuallyOrderedList {
			newItem := newItem.AsValueMap()
			if !slices.ContainsFunc(oldManuallyOrderedList, func(oldItem cty.Value) bool {
				return oldItem.AsValueMap()["order"].AsBigFloat().Cmp(newItem["order"].AsBigFloat()) == 0
			}) {
				newItemOrder, _ := newItem["order"].AsBigFloat().Int64()
				itemToAdd := ObjectRenamingDatabaseManuallyOrderedListItem{
					Name: newItem["name"].AsString(),
					Type: newItem["type"].AsString(),
				}

				// Items can be only added at the end of the list, otherwise invalid operation will be reported.
				if int(newItemOrder) > maxStateOrder {
					finalState = append(finalState, itemToAdd)
					updateChangelog.Added = append(updateChangelog.Added, map[string]any{
						"name": itemToAdd.Name,
						"type": itemToAdd.Type,
					})
				} else {
					invalidOperations = append(invalidOperations, fmt.Errorf("unable to add a new item: %+v, in the middle", itemToAdd))
				}
			}
		}

		if len(invalidOperations) > 0 {
			// Partial is essential in invalid operations because it will prevent invalid state from being saved.
			// It was previously failing the tests because Terraform saves the state automatically despite errors being returned.
			d.Partial(true)
			return diag.FromErr(errors.Join(invalidOperations...))
		} else {
			// Apply the changes. For "normal" implementation instead of sending whole state, single changes should be saved and applied here
			// (places for single actions could be saved based on ObjectRenamingDatabaseInstance.Changelog modifications).
			ObjectRenamingDatabaseInstance.ManuallyOrderedList = finalState
			ObjectRenamingDatabaseInstance.ChangeLog = updateChangelog
		}
	}

	return ReadObjectRenamingListsAndSets(false)(ctx, d, meta)
}

func ReadObjectRenamingListsAndSets(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		list := collections.Map(ObjectRenamingDatabaseInstance.List, func(t ObjectRenamingDatabaseListItem) map[string]any {
			return map[string]any{
				"string": t.String,
				"int":    t.Int,
			}
		})
		if err := d.Set("list", list); err != nil {
			return diag.FromErr(err)
		}

		orderedList := make([]map[string]any, len(ObjectRenamingDatabaseInstance.OrderedList))
		for index, item := range ObjectRenamingDatabaseInstance.OrderedList {
			orderedList[index] = map[string]any{
				"name":  item.Name,
				"order": strconv.Itoa(index),
			}
		}
		if err := d.Set("ordered_list", orderedList); err != nil {
			return diag.FromErr(err)
		}

		itemAdded := len(ObjectRenamingDatabaseInstance.ManuallyOrderedList) > len(d.Get("manually_ordered_list").([]any))
		itemRemoved := len(ObjectRenamingDatabaseInstance.ManuallyOrderedList) < len(d.Get("manually_ordered_list").([]any))
		if withExternalChangesMarking && d.Get("manually_ordered_list") != nil && (itemAdded || itemRemoved) {
			// Detecting external changes by comparing current state with external source
			// Improvements:
			// - When the items' length is the same, try to match items by unique combinations (like name + type in this case).
			// - If items were added externally, see if the item was added at the end (valid operation) or somewhere in the middle (invalid operation).
			// - If items were removed externally, see if the items remained in the same order (valid operation).
			// - Handle cases where multiple external operations were done at once (e.g. added and removed an item).
			return diag.FromErr(errors.New("detected external changes in manually_ordered_list"))
		}

		if d.GetRawState().IsNull() {
			// For the first read, let's "copy-paste" config into state
			if err := resources.SetStateToValuesFromConfig(d, objectRenamingListsAndSetsSchema, []string{"manually_ordered_list"}); err != nil {
				return diag.FromErr(err)
			}
		} else {
			// For later reads, let's put external changes into the state. Because we don't get the information order
			// from the external source, we have to guess it. We do it by matching first with state, later with config (if not found).
			// To correctly find items and their order, you have to match by using fields that uniquely identify a given item (name + type in this case).

			manuallyOrderedList := make([]any, len(ObjectRenamingDatabaseInstance.ManuallyOrderedList))
			for index, item := range ObjectRenamingDatabaseInstance.ManuallyOrderedList {
				var itemOrder int64 = -1

				foundIndex := slices.IndexFunc(d.GetRawState().AsValueMap()["manually_ordered_list"].AsValueSlice(), func(value cty.Value) bool {
					return value.AsValueMap()["name"].AsString() == item.Name && value.AsValueMap()["type"].AsString() == item.Type
				})
				if foundIndex != -1 {
					itemOrder, _ = d.GetRawState().AsValueMap()["manually_ordered_list"].AsValueSlice()[foundIndex].AsValueMap()["order"].AsBigFloat().Int64()
				}

				if foundIndex == -1 && !d.GetRawConfig().IsNull() {
					configFoundIndex := slices.IndexFunc(d.GetRawConfig().AsValueMap()["manually_ordered_list"].AsValueSlice(), func(value cty.Value) bool {
						return value.AsValueMap()["name"].AsString() == item.Name && value.AsValueMap()["type"].AsString() == item.Type
					})
					if configFoundIndex != -1 {
						itemOrder, _ = d.GetRawConfig().AsValueMap()["manually_ordered_list"].AsValueSlice()[configFoundIndex].AsValueMap()["order"].AsBigFloat().Int64()
					}
				}

				manuallyOrderedList[index] = map[string]any{
					"name":  item.Name,
					"type":  item.Type,
					"order": itemOrder,
				}
			}

			if err := d.Set("manually_ordered_list", manuallyOrderedList); err != nil {
				return diag.FromErr(err)
			}
		}

		return nil
	}
}

func DeleteObjectRenamingListsAndSets(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	ObjectRenamingDatabaseInstance.List = nil
	ObjectRenamingDatabaseInstance.OrderedList = nil
	ObjectRenamingDatabaseInstance.ManuallyOrderedList = nil
	d.SetId("")

	return nil
}

func ignoreListOrderAfterFirstApply(parentKey string) schema.SchemaDiffSuppressFunc {
	return func(key string, oldValue string, newValue string, d *schema.ResourceData) bool {
		if strings.HasSuffix(key, ".#") {
			return false
		}

		// Raw state is not null after first apply
		if !d.GetRawState().IsNull() {
			// Parse item index from the key
			keyParts := strings.Split(strings.TrimLeft(key, parentKey+"."), ".")
			if len(keyParts) >= 2 {
				index, err := strconv.Atoi(keyParts[0])
				if err != nil {
					log.Println("[DEBUG] Failed to convert list item index: ", err)
					return false
				}

				newItems := d.GetRawConfig().AsValueMap()[parentKey].AsValueSlice()
				if len(newItems) <= index {
					// item was removed
					return false
				}
				newItemHash := newItems[index].Hash()

				newItemWasAlreadyPresent := false

				// Try to find the same hash in the state; if found, the new item was already present, and it only changed place in the list
				for _, oldItem := range d.GetRawState().AsValueMap()[parentKey].AsValueSlice() {
					// Matching hashes indicate the order changed, but the item stayed in the config, so suppress the change
					if oldItem.Hash() == newItemHash {
						newItemWasAlreadyPresent = true
					}
				}

				oldItemIsStillPresent := false

				// Sizes of config and state may not be the same
				if len(d.GetRawState().AsValueMap()[parentKey].AsValueSlice()) > index {
					// Get the hash of the whole item from state (because it represents old value)
					oldItemHash := d.GetRawState().AsValueMap()[parentKey].AsValueSlice()[index].Hash()

					// Try to find the same hash in the config; if found, the old item still exists, but changed its place in the list
					for _, newItem := range d.GetRawConfig().AsValueMap()[parentKey].AsValueSlice() {
						if newItem.Hash() == oldItemHash {
							oldItemIsStillPresent = true
						}
					}
				} else if newItemWasAlreadyPresent {
					// Happens in cases where there's a new item at the end of the list, but it was already present, so do nothing
					return true
				}

				if newItemWasAlreadyPresent && oldItemIsStillPresent {
					return true
				}
			}
		}

		return false
	}
}

func ignoreOrderAfterFirstApplyWithOrderedList(parentKey string) schema.SchemaDiffSuppressFunc {
	return func(key string, oldValue string, newValue string, d *schema.ResourceData) bool {
		if strings.HasSuffix(key, ".#") {
			return false
		}

		// Raw state is not null after first apply
		if !d.GetRawState().IsNull() {
			// Parse item index from the key
			keyParts := strings.Split(strings.TrimLeft(key, parentKey+"."), ".")
			if len(keyParts) >= 2 {
				index, err := strconv.Atoi(keyParts[0])
				if err != nil {
					log.Println("[DEBUG] Failed to convert list item index: ", err)
					return false
				}

				newItem := d.GetRawConfig().AsValueMap()[parentKey].AsValueSlice()[index]

				newItemOrder := -1
				// The new order value cannot be retrieved because it's not set on config level.
				// There's also no other way (most likely) to get the newly computed order value for a given item,
				// making this approach not possible.
				newItemOrderValue := newItem.AsValueMap()["order"]
				if !newItemOrderValue.IsNull() {
					newItemOrder, _ = strconv.Atoi(newItemOrderValue.AsString())
				}

				_ = newItemOrder
			}
		}

		return false
	}
}

func ignoreOrderAfterFirstApplyWithManuallyOrderedList(parentKey string) schema.SchemaDiffSuppressFunc {
	return func(key string, oldValue string, newValue string, d *schema.ResourceData) bool {
		if strings.HasSuffix(key, ".#") {
			return false
		}

		// Raw state is not null after first apply
		if !d.GetRawState().IsNull() {
			// Parse item index from the key
			keyParts := strings.Split(strings.TrimLeft(key, parentKey+"."), ".")
			if len(keyParts) >= 2 {
				index, err := strconv.Atoi(keyParts[0])
				if err != nil {
					log.Println("[DEBUG] Failed to convert list item index: ", err)
					return false
				}

				newItems := d.GetRawConfig().AsValueMap()[parentKey].AsValueSlice()
				if len(newItems) <= index {
					// item was removed
					return false
				}

				newItem := newItems[index]
				itemWasAlreadyPresent := false
				itemIsStillPresent := false

				var newItemOrder int64
				newItemOrderValue := newItem.AsValueMap()["order"]
				if !newItemOrderValue.IsNull() {
					newItemOrder, _ = newItemOrderValue.AsBigFloat().Int64()
				} else {
					// That's a new item
					return false
				}

				// It was already present, but we need to check the hash
				for _, oldItem := range d.GetRawState().AsValueMap()[parentKey].AsValueSlice() {
					oldItemOrder, _ := oldItem.AsValueMap()["order"].AsBigFloat().Int64()
					if oldItemOrder == newItemOrder {
						if oldItem.Hash() != newItem.Hash() {
							// The item has the same order, but the values in other fields changed (different hash)
							return false
						} else {
							itemWasAlreadyPresent = true
							break
						}
					}
				}

				// Check if a new item is indexable (with new items added at the end, it's not possible to index state value for them, because it doesn't exist yet)
				if len(d.GetRawState().AsValueMap()[parentKey].AsValueSlice()) > index {
					oldItem := d.GetRawState().AsValueMap()[parentKey].AsValueSlice()[index]
					oldItemOrder, _ := oldItem.AsValueMap()["order"].AsBigFloat().Int64()

					// Check if this order is still present
					for _, newItem := range d.GetRawConfig().AsValueMap()[parentKey].AsValueSlice() {
						newItemOrder, _ := newItem.AsValueMap()["order"].AsBigFloat().Int64()
						if oldItemOrder == newItemOrder {
							if oldItem.Hash() != newItem.Hash() {
								// The order is still present, but the values in other fields changed (different hash)
								return false
							} else {
								itemIsStillPresent = true
								break
							}
						}
					}
				}

				if itemWasAlreadyPresent && itemIsStillPresent {
					return true
				}
			}
		}

		return false
	}
}
