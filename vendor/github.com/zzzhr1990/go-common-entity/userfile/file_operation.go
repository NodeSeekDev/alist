package userfile

const (
	// OptionNotLockUserTable lock table
	OptionNotLockUserTable = 1
	// OptionOverrideFile override a file when exists
	OptionOverrideFile = 2
	// OptionInternalForceDirectory INTERNAL USE, determin if need internal call
	OptionInternalForceDirectory = 4
	// OptionForceReadMaster INTERNAL USE, force database read master
	OptionForceReadMaster = 8
	// OptionDisplayDeleteFile INTERNAL USE, list deleted file
	OptionDisplayDeleteFile = 16
	// OptionIgnoreExists ignore file create when exists
	OptionIgnoreExists = 32
	// OptionDestnationIsFile the copy, move destnation is a file, not a directory
	OptionDestnationIsFile = 64
	// OptionForceDeleteFile Delete file not to trash
	OptionForceDeleteFile = 128
	// OptionNotCreateFile not create file
	OptionNotCreateFile = 256
	// OptionIgnoreHeaderFileName ignore the file name
	OptionIgnoreHeaderFileName = 512
	// OptionTerminateWithIllegalFile if enable illegal file cannot be created
	OptionTerminateWithIllegalFile = 1024
)

// NotCreateFile not create file in user directory
func NotCreateFile(op int32) bool {
	return op&OptionNotCreateFile == OptionNotCreateFile
}

// IgnoreHeaderFileName igf
func IgnoreHeaderFileName(op int32) bool {
	return op&OptionIgnoreHeaderFileName == OptionIgnoreHeaderFileName
}

// DestnationIsFile if the file need
func DestnationIsFile(op int32) bool {
	return op&OptionDestnationIsFile == OptionDestnationIsFile
}

// NeedLockTable if the file need
func NeedLockTable(op int32) bool {
	return op&OptionNotLockUserTable != OptionNotLockUserTable
}

// NeedOverrideFile determin the file ov
func NeedOverrideFile(op int32) bool {
	return op&OptionOverrideFile == OptionOverrideFile
}

// NeedForceOperateDirectory determin the file is a directory, continue process
func NeedForceOperateDirectory(op int32) bool {
	return op&OptionInternalForceDirectory == OptionInternalForceDirectory
}

// NeedForceForceReadMaster force read master
func NeedForceForceReadMaster(op int32) bool {
	return op&OptionForceReadMaster == OptionForceReadMaster
}

// NeedDisplayDeleteFile Delete File
func NeedDisplayDeleteFile(op int32) bool {
	return op&OptionDisplayDeleteFile == OptionDisplayDeleteFile
}

// NeedIgnoreExists if File exists, just ignore it
func NeedIgnoreExists(op int32) bool {
	return op&OptionIgnoreExists == OptionIgnoreExists
}

// NeedForceDeleteFile delete file instead of move it to trash
func NeedForceDeleteFile(op int32) bool {
	return op&OptionForceDeleteFile == OptionForceDeleteFile
}
