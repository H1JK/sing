package bufio

import (
	"github.com/sagernet/sing/common/buf"
	N "github.com/sagernet/sing/common/network"
)

func CopyExtendedOnce(dst N.ExtendedWriter, src N.ExtendedReader) (n int64, err error) {
	frontHeadroom := N.CalculateFrontHeadroom(dst)
	rearHeadroom := N.CalculateRearHeadroom(dst)
	bufferSize := N.CalculateMTU(src, dst)
	if bufferSize > 0 {
		bufferSize += frontHeadroom + rearHeadroom
	} else {
		bufferSize = buf.BufferSize
	}
	buffer := buf.NewSize(bufferSize)
	defer buffer.Release()
	readBufferRaw := buffer.Slice()
	readBuffer := buf.With(readBufferRaw[:cap(readBufferRaw)-rearHeadroom])
	readBuffer.Resize(frontHeadroom, 0)
	err = src.ReadBuffer(readBuffer)
	if err != nil {
		return
	}
	dataLen := readBuffer.Len()
	buffer.Resize(readBuffer.Start(), dataLen)
	err = dst.WriteBuffer(buffer)
	n = int64(dataLen)
	return
}
