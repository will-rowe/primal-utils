package gff

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/biogo/biogo/io/featio/gff"
	"github.com/biogo/biogo/seq"
	"github.com/will-rowe/primal-utils/src/rampart"
)

// ExtractFeatures takes a GFF file and returns the features
// TODO: instead of populating a map, return a channel that emits rampart.Gene
func ExtractFeatures(gffFile string) (map[string]rampart.Gene, error) {

	// get the map ready
	mappy := make(map[string]rampart.Gene)

	// open the GFF and get a reader
	fh, err := os.Open(gffFile)
	if err != nil {
		return nil, fmt.Errorf("could not open gff file %s: %s", gffFile, err)
	}
	gffReader := gff.NewReader(bufio.NewReader(fh))
	featureCounter := 0
	for {
		feat, err := gffReader.Read()
		if err == io.EOF {
			break

		} else if err != nil {
			return nil, fmt.Errorf("failed to read feature: %s", err)
		}
		gffFeat, _ := feat.(*gff.Feature)

		// process the feature
		// TODO: currently, only interested in genes?
		switch gffFeat.Feature {
		case "gene":
			featureCounter++

			// TODO: how to handle strands? convert strand into 0 / 1 from - / +
			strand := 0
			if gffFeat.FeatStrand == seq.Plus {
				strand = 1
			}

			// TODO: get a gene name from the attribute, default to the counter if absent
			geneName := strconv.Itoa(featureCounter)
			if len(gffFeat.FeatAttributes) != 0 {

				// update gene name here

			}

			// add the gene to the map
			mappy[geneName] = rampart.Gene{
				Start:  gffFeat.FeatStart,
				End:    gffFeat.FeatEnd,
				Strand: strand,
			}
		default:
			continue
		}
	}

	return mappy, nil
}
