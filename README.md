# primal-utils

> WIP: this is a work in progress...

at the moment, `primal-utils` is a one-trip pony to get a [rampart](https://github.com/artic-network/rampart) config file from the output of [primal scheme](http://primal.zibraproject.org/) - but watch this space for more features

## usage

```bash
primal-utils convert -i example-data/NiV_6_Malaysia/ -g example-data/AJ564622.gff -o config.json
```

## todo

* parse primal scheme bed file and then populate the amplicon field in the rampart config
* currently using the biogo GFF parser, which doesn't handle GFF3, so this needs swapping out
* even better, remove the GFF requirement. GFF isn't provided in the primal scheme output, so need to decide how best to provide this information
* address the `TODO` flags in the code...

## caveats

* assumes that the first entry in the primal scheme multifasta output is the intended reference sequence