package internal

import (
	"io"
	"os"
	"path"

	"github.com/wal-g/tracelog"
	"github.com/wal-g/wal-g/utility"
)

const (
	StreamPrefix           = "stream_"
	StreamBackupNameLength = 23 // len(StreamPrefix) + len(utility.BackupTimeFormat)
)

// TODO : unit tests
// PushStream compresses a stream and push it
func (uploader *Uploader) PushStream(stream io.Reader) (string, error) {
	backupName := StreamPrefix + utility.TimeNowCrossPlatformUTC().Format(utility.BackupTimeFormat)
	dstPath := GetStreamName(backupName, uploader.Compressor.FileExtension())
	err := uploader.PushStreamToDestination(stream, dstPath)

	return backupName, err
}

// TODO : unit tests
// PushStreamToDestination compresses a stream and push it to specifyed destination
func (uploader *Uploader) PushStreamToDestination(stream io.Reader, dstPath string) error {
	if uploader.dataSize != nil {
		stream = NewWithSizeReader(stream, uploader.dataSize)
	}
	compressed := CompressAndEncrypt(stream, uploader.Compressor, ConfigureCrypter())
	err := uploader.Upload(dstPath, compressed)
	tracelog.InfoLogger.Println("FILE PATH:", dstPath)

	return err
}

// FileIsPiped Check if file is piped
func FileIsPiped(stream *os.File) bool {
	stat, _ := stream.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func GetStreamName(backupName string, extension string) string {
	return utility.SanitizePath(path.Join(backupName, "stream.")) + extension
}
