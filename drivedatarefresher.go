package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	invocationCommand := os.Args[1]

	switch invocationCommand {
	case "generatebashfile":
		driveSizeArg := os.Args[2]
		driveSize, err := strconv.ParseUint(driveSizeArg, 10, 64)
		if err != nil {
			fmt.Println("Invalid Command")
			os.Exit(1)
		}

		chunkSizeArg := os.Args[3]
		chunkSize, err := strconv.ParseUint(chunkSizeArg, 10, 64)
		if err != nil {
			fmt.Println("Invalid Command")
			os.Exit(1)
		}

		drivePath := os.Args[4]
		chunkFilePath := os.Args[5]
		bashFilePath := os.Args[6]

		var bashFileBuf []byte

		for i := uint64(0); i < (driveSize / chunkSize); i++ {
			var lineBytes []byte

			lineBytes = []byte("sudo dd if=" + drivePath +
				" of=" + chunkFilePath +
				" bs=" + strconv.FormatUint(chunkSize, 10) +
				" skip=" + strconv.FormatUint(i, 10) +
				" iflag=fullblock" +
				" count=1" +
				" conv=notrunc")

			bashFileBuf = append(bashFileBuf, lineBytes...)
			bashFileBuf = append(bashFileBuf, 0x0a)

			lineBytes = []byte("sudo sync")

			bashFileBuf = append(bashFileBuf, lineBytes...)
			bashFileBuf = append(bashFileBuf, 0x0a)

			lineBytes = []byte("sudo drivedatarefresher xordrivechunk " +
				chunkFilePath)

			bashFileBuf = append(bashFileBuf, lineBytes...)
			bashFileBuf = append(bashFileBuf, 0x0a)

			lineBytes = []byte("sudo sync")

			bashFileBuf = append(bashFileBuf, lineBytes...)
			bashFileBuf = append(bashFileBuf, 0x0a)

			lineBytes = []byte("sudo dd if=" + chunkFilePath +
				" of=" + drivePath +
				" bs=" + strconv.FormatUint(chunkSize, 10) +
				" seek=" + strconv.FormatUint(i, 10) +
				" iflag=fullblock" +
				" count=1" +
				" conv=notrunc")

			bashFileBuf = append(bashFileBuf, lineBytes...)
			bashFileBuf = append(bashFileBuf, 0x0a)

			lineBytes = []byte("sudo sync")

			bashFileBuf = append(bashFileBuf, lineBytes...)
			bashFileBuf = append(bashFileBuf, 0x0a)
		}

		os.WriteFile(bashFilePath, bashFileBuf, 0644)

	case "xordrivechunk":
		chunkFilePath := os.Args[2]

		chunkBuf, err := os.ReadFile(chunkFilePath)
		if err != nil {
			fmt.Println("Unable to read chunk file")
			os.Exit(1)
		}

		for i := 0; i < len(chunkBuf); i++ {
			chunkBuf[i] = chunkBuf[i] ^ 0xa7
		}

		os.WriteFile(chunkFilePath, chunkBuf, 0644)

	default:
		fmt.Println("Invaild command")
		os.Exit(1)
	}
}
