package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/asteris-llc/hammer/hammer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	BuildCmd = &cobra.Command{
		Use:   "build [package...]",
		Short: "build packages",
		Long:  "build all packages by default, unless specific packages are specified",
		Run: func(cmd *cobra.Command, packageNames []string) {
			loader := hammer.NewLoader(viper.GetString("search"))
			loaded, err := loader.Load()
			if err != nil {
				logrus.WithField("error", err).Fatal("could not load packages")
			}

			packages := []*hammer.Package{}
			for _, pkg := range loaded {
				for _, name := range packageNames {
					if pkg.Name == name {
						packages = append(packages, pkg)
						break
					}
				}
			}

			if len(packages) == 0 {
				logrus.Fatal("no packages selected")
			}

			packager := hammer.NewPackager(packages)

			err = packager.EnsureOutputDir(viper.GetString("output"))
			if err != nil {
				logrus.WithField("error", err).Fatal("could not create output directory")
			}

			if !packager.Build() { // Errors are already reported to the user from here
				os.Exit(1)
			}
		},
	}
)