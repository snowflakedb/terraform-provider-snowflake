package gen

import "reflect"

func KindOfT[T any]() string {
	t := reflect.TypeFor[T]()
	return t.Name()
}

func KindOfTPointer[T any]() string {
	return KindOfPointer(KindOfT[T]())
}

func KindOfTSlice[T any]() string {
	return KindOfSlice(KindOfT[T]())
}

func KindOfPointer(kind string) string {
	return "*" + kind
}

func KindOfSlice(kind string) string {
	return "[]" + kind
}
