package architest

type (
	FileFilter   = func(*File) bool
	FileReceiver = func(*File)
	Files        []File
)

func (files Files) Filter(filter FileFilter) Files {
	filteredFiles := make([]File, 0)
	for _, f := range files {
		f := f
		if filter(&f) {
			filteredFiles = append(filteredFiles, f)
		}
	}
	return filteredFiles
}

func (files Files) All(receiver FileReceiver) {
	for _, file := range files {
		file := file
		receiver(&file)
	}
}
