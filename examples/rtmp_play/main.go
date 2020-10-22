package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tyrese/joy4/av"
	"github.com/tyrese/joy4/av/avutil"
	"github.com/tyrese/joy4/format"
	"github.com/tyrese/joy4/format/flv"
	"github.com/tyrese/joy4/format/rtmp"
)

func init() {
	format.RegisterAll()
}

// ./rtmp_play
func main() {
	flag.Usage = func() {
		fmt.Println("./rtmp_play -i='rtmp://localhost:1936/app/publish' | ffplay -")
		fmt.Println()
		flag.PrintDefaults()
	}

	dstUrl := flag.String("i", "rtmp://localhost:1935/live/test", "input url")
	debug := flag.Bool("v", false, "verbose")
	quic := flag.Bool("q", false, "rtmp over quic")
	dumpFile := flag.String("f", "stdout", "dump to file")
	flag.Parse()
	rtmp.Debug = *debug
	rtmp.UseQuic = *quic
	conn, _ := rtmp.Dial(*dstUrl)

	var muxer av.Muxer
	if *dumpFile == "stdout" {
		muxer = flv.NewMuxer(os.Stdout)
	} else {
		fmt.Println("dump to file " + *dumpFile)
		os.Remove(*dumpFile)
		file, err := os.Create(*dumpFile)
		if err != nil {
			fmt.Println(err.Error())
		}

		muxer = flv.NewMuxer(file)
	}

	avutil.CopyFile(muxer, conn)

	conn.Close()
}
