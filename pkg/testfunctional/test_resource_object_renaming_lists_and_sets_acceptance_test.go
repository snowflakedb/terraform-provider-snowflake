package testfunctional_test

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SdkV2Functional_BasicListFlow(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "111", "int": 111},
					{"string": "222", "int": 222},
					{"string": "333", "int": 333},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"string": "111", "int": "111"},
						{"string": "222", "int": "222"},
						{"string": "333", "int": "333"},
					}),
				),
			},
			// Remove, shift, and add one item (in the middle)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Removed: []map[string]any{
								{"string": "111", "int": 111},
							},
							Added: []map[string]any{
								{"string": "444", "int": 444},
							},
						}),
					},
				},
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "222", "int": 222},
					{"string": "444", "int": 444},
					{"string": "333", "int": 333},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"string": "222", "int": "222"},
						{"string": "444", "int": "444"},
						{"string": "333", "int": "333"},
					}),
				),
			},
			// Remove, shift, and add one item (at the end)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Removed: []map[string]any{
								{"string": "222", "int": 222},
							},
							Added: []map[string]any{
								{"string": "111", "int": 111},
							},
						}),
					},
				},
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "444", "int": 444},
					{"string": "333", "int": 333},
					{"string": "111", "int": 111},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"string": "444", "int": "444"},
						{"string": "333", "int": "333"},
						{"string": "111", "int": "111"},
					}),
				),
			},
			// Remove, shift, and add one item (at the beginning)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Removed: []map[string]any{
								{"string": "111", "int": 111},
							},
							Added: []map[string]any{
								{"string": "222", "int": 222},
							},
						}),
					},
				},
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "222", "int": 222},
					{"string": "333", "int": 333},
					{"string": "444", "int": 444},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"string": "222", "int": "222"},
						{"string": "333", "int": "333"},
						{"string": "444", "int": "444"},
					}),
				),
			},
			// Reorder items and add one
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Added: []map[string]any{
								{"string": "555", "int": 555},
							},
						}),
					},
				},
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "444", "int": 444},
					{"string": "555", "int": 555},
					{"string": "333", "int": 333},
					{"string": "222", "int": 222},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"string": "444", "int": "444"},
						{"string": "555", "int": "555"},
						{"string": "333", "int": "333"},
						{"string": "222", "int": "222"},
					}),
				),
			},
			// Replace all items
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Removed: []map[string]any{
								{"string": "333", "int": 333},
								{"string": "444", "int": 444},
								{"string": "222", "int": 222},
								{"string": "555", "int": 555},
							},
							Added: []map[string]any{
								{"string": "1111", "int": 1111},
								{"string": "2222", "int": 2222},
								{"string": "3333", "int": 3333},
								{"string": "4444", "int": 4444},
							},
						}),
					},
				},
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "1111", "int": 1111},
					{"string": "2222", "int": 2222},
					{"string": "3333", "int": 3333},
					{"string": "4444", "int": 4444},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"string": "1111", "int": "1111"},
						{"string": "2222", "int": "2222"},
						{"string": "3333", "int": "3333"},
						{"string": "4444", "int": "4444"},
					}),
				),
			},
			// Remove a few items
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Removed: []map[string]any{
								{"string": "3333", "int": 3333},
								{"string": "4444", "int": 4444},
							},
						}),
					},
				},
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "1111", "int": 1111},
					{"string": "2222", "int": 2222},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"string": "1111", "int": "1111"},
						{"string": "2222", "int": "2222"},
					}),
				),
			},
			// Remove all items
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Removed: []map[string]any{
								{"string": "1111", "int": 1111},
								{"string": "2222", "int": 2222},
							},
						}),
					},
				},
				Config: objectRenamingConfigList([]map[string]any{}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.#", "0"),
				),
			},
			// Add few items
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Added: []map[string]any{
								{"string": "1111", "int": 1111},
								{"string": "2222", "int": 2222},
							},
						}),
					},
				},
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "1111", "int": 1111},
					{"string": "2222", "int": 2222},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"string": "1111", "int": "1111"},
						{"string": "2222", "int": "2222"},
					}),
				),
			},
			// External changes: add item
			{
				PreConfig: func() {
					testfunctional.ObjectRenamingDatabaseInstance.List = append(testfunctional.ObjectRenamingDatabaseInstance.List, testfunctional.ObjectRenamingDatabaseListItem{
						String: "3333",
						Int:    3333,
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Removed: []map[string]any{
								{"string": "3333", "int": 3333},
							},
						}),
					},
				},
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "1111", "int": 1111},
					{"string": "2222", "int": 2222},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"string": "1111", "int": "1111"},
						{"string": "2222", "int": "2222"},
					}),
				),
			},
			// External changes: removed item
			{
				PreConfig: func() {
					testfunctional.ObjectRenamingDatabaseInstance.List = testfunctional.ObjectRenamingDatabaseInstance.List[:1]
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Added: []map[string]any{
								{"string": "2222", "int": 2222},
							},
						}),
					},
				},
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "1111", "int": 1111},
					{"string": "2222", "int": 2222},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"string": "1111", "int": "1111"},
						{"string": "2222", "int": "2222"},
					}),
				),
			},
			// External changes: change item
			{
				PreConfig: func() {
					testfunctional.ObjectRenamingDatabaseInstance.List[1].String = "1010"
					testfunctional.ObjectRenamingDatabaseInstance.List[1].Int = 1010
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Removed: []map[string]any{
								{"string": "1010", "int": 1010},
							},
							Added: []map[string]any{
								{"string": "2222", "int": 2222},
							},
						}),
					},
				},
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "1111", "int": 1111},
					{"string": "2222", "int": 2222},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"string": "1111", "int": "1111"},
						{"string": "2222", "int": "2222"},
					}),
				),
			},
			// Add an item that is identical to another one (currently, failing because hash duplicates are not handled)
			// {
			//	Config: objectRenamingConfigList([]map[string]any{
			//		{"string": "222", "int": 222},
			//		{"string": "222", "int": 222},
			//	}),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		assert.HasListItemsOrderIndependent("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
			//			{"string": "222", "int": "222"},
			//			{"string": "222", "int": "222"},
			//		}),
			//	),
			// },
		},
	})
}

// This test researches the possibility of performing update instead of remove + add item
func TestAcc_SdkV2Functional_ListNameUpdate(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	t.Skip("Currently failing, will be investigated soon.")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"name": "column1", "string": "111", "int": 111},
					{"name": "column2", "string": "222", "int": 222},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"name": "column1", "string": "111", "int": "111"},
						{"name": "column2", "string": "222", "int": "222"},
					}),
				),
			},
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"name": "column2", "string": "222", "int": 222},
					{"name": "column1", "string": "111", "int": 111},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"name": "column2", "string": "222", "int": "222"},
						{"name": "column1", "string": "111", "int": "111"},
					}),
				),
			},
			// It's hard to handle reorder + rename with this approach,
			// because without any additional metadata, we cannot identify a given list item.
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"name": "column3", "string": "222", "int": 222},
					{"name": "column1", "string": "111", "int": 111},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"name": "column3", "string": "222", "int": "222"},
						{"name": "column1", "string": "111", "int": "111"},
					}),
				),
			},
		},
	})
}

func TestAcc_SdkV2Functional_ListsWithDuplicatedItems(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	// Fails, because the SuppressDiffFunc works on the hash of individual items.
	// To correctly suppress such changes, the number of repeated hashes should be counted.
	t.Skip("Currently failing, because duplicated hashes are not supported.")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "111", "int": 111},
					{"string": "222", "int": 222},
					{"string": "333", "int": 333},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"string": "111", "int": "111"},
						{"string": "222", "int": "222"},
						{"string": "333", "int": "333"},
					}),
				),
			},
			// Introduce duplicates (it would be enough just to introduce only one to break the approach assumptions)
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "111", "int": 111},
					{"string": "111", "int": 111},
					{"string": "222", "int": 222},
					{"string": "222", "int": 222},
					{"string": "333", "int": 333},
					{"string": "333", "int": 333},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.ContainsExactlyInAnyOrder("snowflake_test_resource_object_renaming.test", "list", []map[string]string{
						{"string": "111", "int": "111"},
						{"string": "111", "int": "111"},
						{"string": "222", "int": "222"},
						{"string": "222", "int": "222"},
						{"string": "333", "int": "333"},
						{"string": "333", "int": "333"},
					}),
				),
			},
		},
	})
}

func objectRenamingConfigList(listItems []map[string]any) string {
	generateListItem := func(s string, i int) string {
		return fmt.Sprintf(`
	list {
		string = "%[1]s"
		int = %[2]d
	}
`, s, i)
	}

	generatedListItems := ""
	for _, item := range listItems {
		generatedListItems += generateListItem(item["string"].(string), item["int"].(int))
	}

	return fmt.Sprintf(`

	resource "snowflake_test_resource_object_renaming" "test" {
		provider = "%[2]s"

		%[1]s
	}

`, generatedListItems, SdkV2FunctionalTestsProviderName)
}

type objectRenamingPlanCheck func(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse)

func (fn objectRenamingPlanCheck) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	fn(ctx, req, resp)
}

func assertObjectRenamingDatabaseChangelogAndClearIt(changelog testfunctional.ObjectRenamingDatabaseChangelog) plancheck.PlanCheck {
	return objectRenamingPlanCheck(func(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
		if !reflect.DeepEqual(testfunctional.ObjectRenamingDatabaseInstance.ChangeLog, changelog) {
			resp.Error = fmt.Errorf("expected %+v changelog for this step, but got: %+v", changelog, testfunctional.ObjectRenamingDatabaseInstance.ChangeLog)
		}
		testfunctional.ObjectRenamingDatabaseInstance.ChangeLog.Added = nil
		testfunctional.ObjectRenamingDatabaseInstance.ChangeLog.Removed = nil
		testfunctional.ObjectRenamingDatabaseInstance.ChangeLog.Changed = nil
	})
}

func TestAcc_SdkV2Functional_SupportedActions(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "TEXT", "order": 10},
					{"name": "nameTwo", "type": "STRING", "order": 20},
					{"name": "nameThree", "type": "NUMBER", "order": 30},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOne"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "TEXT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameTwo"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "20"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.name", "nameThree"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.type", "NUMBER"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.order", "30"),
				),
			},
			// Drop item (any position)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Removed: []map[string]any{
								{"name": "nameTwo", "type": "STRING"},
							},
						}),
					},
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "TEXT", "order": 10},
					{"name": "nameThree", "type": "NUMBER", "order": 30},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOne"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "TEXT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameThree"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "NUMBER"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "30"),
				),
			},
			// Add item (at the end)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Added: []map[string]any{
								{"name": "nameFour", "type": "INT"},
							},
						}),
					},
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "TEXT", "order": 10},
					{"name": "nameThree", "type": "NUMBER", "order": 30},
					{"name": "nameFour", "type": "INT", "order": 40},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOne"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "TEXT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameThree"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "NUMBER"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "30"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.name", "nameFour"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.type", "INT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.order", "40"),
				),
			},
			// Rename item / Change item type (any position; compatible type change)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Changed: []testfunctional.ObjectRenamingDatabaseChangelogChange{
								{
									Before: map[string]any{"name": "nameOne", "type": "TEXT"},
									After:  map[string]any{"name": "nameOneV2", "type": "STRING"},
								},
								{
									Before: map[string]any{"name": "nameThree", "type": "NUMBER"},
									After:  map[string]any{"name": "nameThreeV2", "type": "INT"},
								},
								{
									Before: map[string]any{"name": "nameFour", "type": "INT"},
									After:  map[string]any{"name": "nameFourV2", "type": "NUMBER"},
								},
							},
						}),
					},
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOneV2", "type": "STRING", "order": 10},
					{"name": "nameThreeV2", "type": "INT", "order": 30},
					{"name": "nameFourV2", "type": "NUMBER", "order": 40},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOneV2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameThreeV2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "INT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "30"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.name", "nameFourV2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.type", "NUMBER"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.order", "40"),
				),
			},
			// Reorder items in the configuration
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionNoop),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{}),
					},
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameFourV2", "type": "NUMBER", "order": 40},
					{"name": "nameThreeV2", "type": "INT", "order": 30},
					{"name": "nameOneV2", "type": "STRING", "order": 10},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOneV2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameThreeV2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "INT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "30"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.name", "nameFourV2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.type", "NUMBER"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.order", "40"),
				),
			},
			// (after reorder) Drop item (any position)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Removed: []map[string]any{
								{"name": "nameThreeV2", "type": "INT"},
							},
						}),
					},
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameFourV2", "type": "NUMBER", "order": 40},
					{"name": "nameOneV2", "type": "STRING", "order": 10},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOneV2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameFourV2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "NUMBER"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "40"),
				),
			},
			// (after reorder) Add item (at the end)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Added: []map[string]any{
								{"name": "nameFive", "type": "INT"},
							},
						}),
					},
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameFive", "type": "INT", "order": 50},
					{"name": "nameFourV2", "type": "NUMBER", "order": 40},
					{"name": "nameOneV2", "type": "STRING", "order": 10},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOneV2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameFourV2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "NUMBER"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "40"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.name", "nameFive"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.type", "INT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.order", "50"),
				),
			},
			// (after reorder) Rename item / Change item type (any position; compatible type change)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{
							Changed: []testfunctional.ObjectRenamingDatabaseChangelogChange{
								{
									Before: map[string]any{"name": "nameOneV2", "type": "STRING"},
									After:  map[string]any{"name": "nameOneV10", "type": "TEXT"},
								},
								{
									Before: map[string]any{"name": "nameFourV2", "type": "NUMBER"},
									After:  map[string]any{"name": "nameFourV10", "type": "INT"},
								},
								{
									Before: map[string]any{"name": "nameFive", "type": "INT"},
									After:  map[string]any{"name": "nameFiveV10", "type": "NUMBER"},
								},
							},
						}),
					},
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameFiveV10", "type": "NUMBER", "order": 50},
					{"name": "nameFourV10", "type": "INT", "order": 40},
					{"name": "nameOneV10", "type": "TEXT", "order": 10},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOneV10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "TEXT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameFourV10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "INT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "40"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.name", "nameFiveV10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.type", "NUMBER"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.order", "50"),
				),
			},
		},
	})
}

func TestAcc_SdkV2Functional_UnsupportedActions_AddItemsNotAtTheEnd(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameTwo", "type": "STRING", "order": 20},
					{"name": "nameOne", "type": "TEXT", "order": 10},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOne"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "TEXT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameTwo"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "20"),
				),
			},
			// Add item (in the middle)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(
							"snowflake_test_resource_object_renaming.test",
							"manually_ordered_list",
							tfjson.ActionUpdate,
							sdk.String("[map[name:nameOne order:10 type:TEXT] map[name:nameTwo order:20 type:STRING]]"),
							sdk.String("[map[name:atTheBeginning order:15 type:TEXT] map[name:nameTwo order:20 type:STRING] map[name:inTheMiddle order:17 type:INT] map[name:nameTwo order:20 type:STRING]]"),
						),
					},
					// PostChecks don't apply when the expected error is set
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "atTheBeginning", "type": "TEXT", "order": 15},
					{"name": "nameOne", "type": "TEXT", "order": 10},
					{"name": "inTheMiddle", "type": "INT", "order": 17},
					{"name": "nameTwo", "type": "STRING", "order": 20},
				}),
				ExpectError: regexp.MustCompile("unable to add a new item: \\{Name:atTheBeginning Type:TEXT}, in the middle\nunable to add a new item: \\{Name:inTheMiddle Type:INT}, in the middle"),
			},
			// Try to go back to the original state (with flipped items in config)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{}),
					},
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "TEXT", "order": 10},
					{"name": "nameTwo", "type": "STRING", "order": 20},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOne"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "TEXT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameTwo"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "20"),
				),
			},
		},
	})
}

func TestAcc_SdkV2Functional_UnsupportedActions_ChangeItemTypeToIncompatibleOne(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameTwo", "type": "STRING", "order": 20},
					{"name": "nameOne", "type": "TEXT", "order": 10},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOne"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "TEXT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameTwo"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "20"),
				),
			},
			// Change item type (incompatible change)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(
							"snowflake_test_resource_object_renaming.test",
							"manually_ordered_list",
							tfjson.ActionUpdate,
							sdk.String("[map[name:nameOne order:10 type:TEXT] map[name:nameTwo order:20 type:STRING]]"),
							sdk.String("[map[name:nameOne order:10 type:NUMBER] map[name:nameTwo order:20 type:STRING]]"),
						),
					},
					// PostChecks don't apply when the expected error is set
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "NUMBER", "order": 10},
					{"name": "nameTwo", "type": "STRING", "order": 20},
				}),
				ExpectError: regexp.MustCompile("unable to change item type from TEXT to NUMBER"),
			},
			// Try to go back to the original state (with flipped items in config)
			{
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "TEXT", "order": 10},
					{"name": "nameTwo", "type": "STRING", "order": 20},
				}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionNoop),
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{}),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOne"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "TEXT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameTwo"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "20"),
				),
			},
		},
	})
}

func TestAcc_SdkV2Functional_UnsupportedActions_ExternalChange_AddNewItem(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "TEXT", "order": 10},
					{"name": "nameTwo", "type": "STRING", "order": 20},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOne"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "TEXT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameTwo"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "20"),
				),
			},
			// Add one item externally
			{
				PreConfig: func() {
					testfunctional.ObjectRenamingDatabaseInstance.ManuallyOrderedList = append(testfunctional.ObjectRenamingDatabaseInstance.ManuallyOrderedList, testfunctional.ObjectRenamingDatabaseManuallyOrderedListItem{
						Name: "externalItem",
						Type: "INT",
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(
							"snowflake_test_resource_object_renaming.test",
							"manually_ordered_list",
							tfjson.ActionUpdate,
							sdk.String("[map[name:nameOne order:10 type:TEXT] map[name:nameTwo order:20 type:STRING]]"),
							sdk.String("[map[name:nameOne order:10 type:NUMBER] map[name:nameTwo order:20 type:STRING] map[name:externalItem order:-1 type:INT]]"),
						),
					},
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "NUMBER", "order": 10},
					{"name": "nameTwo", "type": "STRING", "order": 20},
				}),
				ExpectError: regexp.MustCompile("detected external changes in manually_ordered_list"),
			},
			// Try to go back to the original state (after external correction)
			{
				PreConfig: func() {
					testfunctional.ObjectRenamingDatabaseInstance.ManuallyOrderedList = testfunctional.ObjectRenamingDatabaseInstance.ManuallyOrderedList[:len(testfunctional.ObjectRenamingDatabaseInstance.ManuallyOrderedList)-1]
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{}),
					},
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "TEXT", "order": 10},
					{"name": "nameTwo", "type": "STRING", "order": 20},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOne"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "TEXT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameTwo"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "20"),
				),
			},
		},
	})
}

func TestAcc_SdkV2Functional_UnsupportedActions_ExternalChange_RemoveItem(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "TEXT", "order": 10},
					{"name": "nameThree", "type": "INT", "order": 30},
					{"name": "nameTwo", "type": "STRING", "order": 20},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.#", "3"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOne"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "TEXT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameTwo"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "20"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.name", "nameThree"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.type", "INT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.order", "30"),
				),
			},
			// Remove one item externally
			{
				PreConfig: func() {
					// Remove middle item
					testfunctional.ObjectRenamingDatabaseInstance.ManuallyOrderedList = append(testfunctional.ObjectRenamingDatabaseInstance.ManuallyOrderedList[:1], testfunctional.ObjectRenamingDatabaseInstance.ManuallyOrderedList[2])
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(
							"snowflake_test_resource_object_renaming.test",
							"manually_ordered_list",
							tfjson.ActionUpdate,
							sdk.String("[map[name:nameOne order:10 type:NUMBER] map[name:nameTwo order:20 type:STRING] map[name:nameThree order:30 type:INT]]"),
							sdk.String("[map[name:nameOne order:10 type:NUMBER] map[name:nameThree order:30 type:INT]]"),
						),
					},
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "TEXT", "order": 10},
					{"name": "nameThree", "type": "INT", "order": 30},
					{"name": "nameTwo", "type": "STRING", "order": 20},
				}),
				ExpectError: regexp.MustCompile("detected external changes in manually_ordered_list"),
			},
			// Try to go back to the original state (after external correction)
			{
				PreConfig: func() {
					// Bring the middle item back
					testfunctional.ObjectRenamingDatabaseInstance.ManuallyOrderedList = []testfunctional.ObjectRenamingDatabaseManuallyOrderedListItem{
						testfunctional.ObjectRenamingDatabaseInstance.ManuallyOrderedList[0],
						{
							Name: "nameTwo",
							Type: "STRING",
						},
						testfunctional.ObjectRenamingDatabaseInstance.ManuallyOrderedList[1],
					}
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						assertObjectRenamingDatabaseChangelogAndClearIt(testfunctional.ObjectRenamingDatabaseChangelog{}),
					},
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "TEXT", "order": 10},
					{"name": "nameThree", "type": "INT", "order": 30},
					{"name": "nameTwo", "type": "STRING", "order": 20},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.#", "3"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOne"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "TEXT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameTwo"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "20"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.name", "nameThree"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.type", "INT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.order", "30"),
				),
			},
		},
	})
}

func TestAcc_SdkV2Functional_UnsupportedActions_ChangingTheOrderOfItem(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "TEXT", "order": 10},
					{"name": "nameThree", "type": "INT", "order": 30},
					{"name": "nameTwo", "type": "STRING", "order": 20},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.#", "3"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOne"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "TEXT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameTwo"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "20"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.name", "nameThree"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.type", "INT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.order", "30"),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionUpdate),
					},
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "TEXT", "order": 15},
					{"name": "nameThree", "type": "INT", "order": 35},
					{"name": "nameTwo", "type": "STRING", "order": 25},
				}),
				ExpectError: regexp.MustCompile("unable to add a new item: \\{Name:nameOne Type:TEXT}, in the middle\nunable to add a new item: \\{Name:nameTwo Type:STRING}, in the middle"),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_test_resource_object_renaming.test", plancheck.ResourceActionNoop),
					},
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "nameOne", "type": "TEXT", "order": 10},
					{"name": "nameThree", "type": "INT", "order": 30},
					{"name": "nameTwo", "type": "STRING", "order": 20},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.#", "3"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.name", "nameOne"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.type", "TEXT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.0.order", "10"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.name", "nameTwo"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.type", "STRING"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.1.order", "20"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.name", "nameThree"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.type", "INT"),
					resource.TestCheckResourceAttr("snowflake_test_resource_object_renaming.test", "manually_ordered_list.2.order", "30"),
				),
			},
		},
	})
}

func objectRenamingConfigManuallyOrderedList(listItems []map[string]any) string {
	generateListItem := func(name string, itemType string, order int) string {
		return fmt.Sprintf(`
manually_ordered_list {
	name = "%s"
	type = "%s"
	order = %d
}
`, name, itemType, order)
	}

	generatedListItems := ""
	for _, item := range listItems {
		generatedListItems += generateListItem(item["name"].(string), item["type"].(string), item["order"].(int))
	}

	return fmt.Sprintf(`
resource "snowflake_test_resource_object_renaming" "test" {
	provider = "%[2]s"

	%[1]s
}
`, generatedListItems, SdkV2FunctionalTestsProviderName)
}
