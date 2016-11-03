package cmds

import (
	"flag"
	"fmt"

	"github.com/tiago4orion/enzo/disk/mbr"
)

var (
	flags                                   *flag.FlagSet
	flagHelp, flagCreate, flagUpdate        *bool
	flagAddpart, flagDelpart, flagStartsect *int
	flagBootcode, flagLastsect              *string
)

func init() {
	flags = flag.NewFlagSet("mbr", flag.ContinueOnError)
	flagHelp = flags.Bool("help", false, "Show this help")
	flagCreate = flags.Bool("create", false, "Create new MBR")
	flagUpdate = flags.Bool("update", false, "Update MBR")
	flagAddpart = flag.Int("add-part", 0, "Add partition")
	flagDelpart = flag.Int("del-part", 0, "Delete partition")
	flagStartsect = flag.Int("start-sect", 0, "start sector")
	flagLastsect = flag.String("last-sect", "", "last sector (modififers +K, +M, +G works)")
	flagBootcode = flags.String("bootcode", "", "Bootsector binary code")
}

func addPart(disk string, partnumber int, startsect int, lastsect string) error {
	mbrdata, err := mbr.FromFile(disk)

	if err != nil {
		return err
	}

	p := mbr.NewEmptyPartition()
	mbrdata.SetPart(partnumber, p)
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
		if *flagAddpart <= 0 || *flagDelpart == 0 {
			return fmt.Errorf("-update requires flag -add-part or --del-part")
		}

		if *flagAddpart != 0 {
			partNumber := *flagAddpart

			if *flagStartsect == -1 {
				return fmt.Errorf("-add-part requires -start-sect")
			}

			if *flagLastsect == "" {
				return fmt.Errorf("-add-part requires -last-sect")
			}

			return addPart(disks[0], partNumber, *flagStartsect, *flagLastsect)
		}

		return fmt.Errorf("-del-part not implemented")
	}

	for _, disk := range disks {
		err := mbr.Info(disk)

		if err != nil {
			return err
		}
	}

	return nil
}
