package archtest

type FileFilter = func(*File) bool
type FileReceiver = func(*File)

func FilterFiles(files []File, filter FileFilter) []File {
	filteredFiles := make([]File, 0)
	for _, f := range files {
		if filter(&f) {
			filteredFiles = append(filteredFiles, f)
		}
	}
	return filteredFiles
}

func IterateFiles(files []File, receiver FileReceiver) {
	for _, file := range files {
		receiver(&file)
	}
}
