# msa2ckmers

## Aim:

This is a hacky python script for testing some primer design ideas.

The aim is to go from `reference genomes` > `graph` > `linear regions` > `k-mers`

## Approach:

* start easy with the MSA of the EBOV reference genomes from ARTIC

* convert MSA > GFA (using `msa2gfa.py`)

* extract all segments from the GFA

* use the GFA links to determine the linear regions

* filter the extracted segments to keep only the linear segments


## Caveats of linear segment finder:

* no tests

* no consideration given to segment orientation in paths

* no k-mer size strategy

* using LN tag in GFA to indicate MSA position of segment