#!/usr/bin/env python3
import sys
import re

kSize = 21

def main():
    paths = []
    segments = {}
    msaPositions = {}

    with open("ebov-ref-genomes.gfa", 'r') as fh:
        for line in fh:

            # add all the path lines to a list
            if line[0] == 'P':
                paths.append(line.rstrip().split()[2])

            # add all segments to a dict, with id as key and seq as value
            elif line[0] == 'S':
                segLine = line.rstrip().split()
                segments[segLine[1]] = segLine[2]
                pos = re.sub('[LN:i:]', '', segLine[3])
                msaPositions[segLine[1]] = pos

    numPaths = len(paths)

    print("number of paths in graph: {}" .format(numPaths))
    print("number of segments in graph: {}" .format(len(segments)))

    # pathSegCount records the number of times a segment appears in the paths
    # the key is the seg ID, the value is the count
    pathSegCount = {}
    for path in paths:
        for seg in path.split(','):
            seg = re.sub('[+]', '', seg)
            if '-' in seg:
                sys.exit('encountered multiple orientations in GFA')
            if seg in pathSegCount:
                pathSegCount[seg] += 1
            else:
                pathSegCount[seg] = 1

    for seg, count in pathSegCount.items():
        if count != numPaths:
            del segments[seg]
        
    print("number of conserved linear segments in graph: {}" .format(len(segments)))

    print("MSA position\tconserved sequence")
    for seg in segments:
        print("{}\t{}" .format(msaPositions[seg], segments[seg]))





if __name__ == '__main__':
    main()