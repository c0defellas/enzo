package cmds

import (
	"flag"
	"fmt"

	"github.com/tiago4orion/enzo/disk/mbr"
)

var (
	flags                            *flag.FlagSet
	flagHelp, flagCreate, flagUpdate *bool
	flagAddpart, flagStartsect       *int
	flagBootcode, flagLastsect       *string
)

func init() {
	flags = flag.NewFlagSet("mbr", flag.ContinueOnError)
	flagHelp = flags.Bool("help", false, "Show this help")
	flagCreate = flags.Bool("create", false, "Create new MBR")
	flagUpdate = flags.Bool("update", false, "Update MBR")
	flagAddpart = flag.Int("add-part", 0, "Add partition")
	flagStartsect = flag.Int("start-sect", 0, "start sector")
	flagLastsect = flag.String("last-sect", "", "last sector (modififers +K, +M, +G works)")
	flagBootcode = flags.String("bootcode", "", "Bootsector binary code")
}

func updateMBR(partnumber int, startsect int, lastsect string) error {
	return nil
}

func MBR(args []string) error {
	flags.Parse(args[1:])

	if *flagHelp {
		flags.PrintDefaults()
		return nil
	}

	disks := flags.Args()

	if len(disks) != 1 {
		return fmt.Errorf("Require one device file")
	}

	if *flagCreate {
		if *flagUpdate {
			return fmt.Errorf("-create conflicts with -update")
		}

		return mbr.Create(disks[0], *flagBootcode)
	}

	if *flagUpdate {
		if *flagAddpart <= 0 {
			return fmt.Errorf("-update requires flag -add-part")
		}

		partNumber := *flagAddpart

		if *flagStartsect == -1 {
			return fmt.Errorf("-add-part requires -start-sect")
		}

		if *flagLastsect == "" {
			return fmt.Errorf("-add-part requires -last-sect")
		}

		return updateMBR(partNumber, *flagStartsect, *flagLastsect)
	}

	for _, disk := range disks {
		err := mbr.Info(disk)

		if err != nil {
			return err
		}
	}

	return nil
}
