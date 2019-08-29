/*
Copyright © 2019 Will Rowe <w.p.m.rowe@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/io/featio/bed"
	"github.com/biogo/biogo/io/seqio"
	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"
	"github.com/spf13/cobra"

	"github.com/will-rowe/primal-utils/src/gff"
	"github.com/will-rowe/primal-utils/src/misc"
	"github.com/will-rowe/primal-utils/src/rampart"
)

// the command line arguments
var (
	primalDir      *string
	gffFile        *string
	outFile        *string
	defaultOutFile = "./rampartConfig-" + string(time.Now().Format("20060102150405")+".json")
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Massage the output of Primal Scheme into a RAMPART config",
	Long: `Massage the output of Primal Scheme into a RAMPART config

	`,
	Run: func(cmd *cobra.Command, args []string) {
		runConvert()
	},
}

func init() {
	primalDir = convertCmd.Flags().StringP("primalDir", "i", "", "directory containing output from Primal Scheme")
	gffFile = convertCmd.Flags().StringP("gff", "g", "", "GFF3 file for reference sequence")
	outFile = convertCmd.Flags().StringP("configFile", "o", defaultOutFile, "name for the RAMPART config")
	convertCmd.MarkFlagRequired("primalDir")
	convertCmd.Flags().SortFlags = false
	rootCmd.AddCommand(convertCmd)
}

func runConvert() {
	log.SetOutput(os.Stdout)
	log.Println("converting Primal Scheme output to a RAMPART config...")
	log.Printf("\tPrimal Scheme directory: %s", *primalDir)

	// check the Primal Scheme directory exists
	misc.ErrorCheck(misc.CheckDir(*primalDir))

	// make sure there is a bed file and multifasta in there
	bedFiles, err := misc.CollectFiles(*primalDir, "bed", false)
	misc.ErrorCheck(err)
	if len(bedFiles) == 0 {
		misc.ErrorCheck(fmt.Errorf("no bed files found in supplied Primal Scheme directory"))
	} else if len(bedFiles) > 1 {
		misc.ErrorCheck(fmt.Errorf("multiple bed files found in supplied Primal Scheme directory, not sure which to use"))
	} else {
		log.Printf("\tfound bed file: %s", bedFiles[0])
	}
	mfastaFiles, err := misc.CollectFiles(*primalDir, "fasta", false)
	misc.ErrorCheck(err)
	if len(mfastaFiles) == 0 {
		misc.ErrorCheck(fmt.Errorf("no multifasta files found in supplied Primal Scheme directory"))
	} else if len(bedFiles) > 1 {
		misc.ErrorCheck(fmt.Errorf("multiple multifasta files found in supplied Primal Scheme directory, not sure which to use"))
	} else {
		log.Printf("\tfound multifasta file: %s", mfastaFiles[0])
	}

	// check the supplied GFF or try getting one
	if *gffFile != "" {
		misc.ErrorCheck(misc.CheckFile(*gffFile))
		log.Printf("\tGFF file: %s", *gffFile)
	} else {
		log.Printf("no gff provided, in future will try finding one online")
		os.Exit(1)
	}

	// get the reference sequence from the multifasta NOTE: this assume the top entry is the reference
	var r *fasta.Reader
	t := linear.NewSeq("", nil, alphabet.DNA)
	if mfasta, err := os.Open(mfastaFiles[0]); err != nil {
		log.Fatalf("failed to open %q: %v", mfastaFiles[0], err)
		os.Exit(1)
	} else {
		defer mfasta.Close()
		r = fasta.NewReader(mfasta, t)
	}

	// create a reference for RAMPART
	newRef := rampart.NewReference()
	sc := seqio.NewScanner(r)
	sc.Next()
	s := sc.Seq()
	newRef.Label = s.Name()
	newRef.Accession = s.Name()
	newRef.Length = s.Len()
	// this bit is rubbish - get a better way to get the seq from a seq.Sequence
	seq := make([]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		ql := s.At(i)
		seq[i] = (string(ql.L))
	}
	newRef.Sequence = strings.Join(seq, "")

	// get the GFF info for the reference
	featureMap, err := gff.ExtractFeatures(*gffFile)
	misc.ErrorCheck(err)
	newRef.Genes = featureMap

	// get the BED info for the reference
	amplicons := [][]int{}
	fh, err := os.Open(bedFiles[0])
	misc.ErrorCheck(err)
	bedReader, err := bed.NewReader(bufio.NewReader(fh), 4)
	misc.ErrorCheck(err)
	for {
		entry, err := bedReader.Read()
		if err == io.EOF {
			break
		}
		misc.ErrorCheck(err)
		amplicons = append(amplicons, []int{entry.Start(), entry.End()})
	}
	newRef.Amplicons = amplicons

	// create the config struct with the newly created reference
	newConf := rampart.NewConfig(newRef)

	// write to disk
	misc.ErrorCheck(newConf.WriteConfig(*outFile))
	log.Printf("config written to: %s", *outFile)
	log.Println("donzo.")
}
