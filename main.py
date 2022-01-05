# (c) Mateusz Slaski for OVHcloud

import sys
import os
import enum
import pandas as pd

class Ages(enum.Enum):
    oldest = 0
    youngest = 1

class RankEnds(enum.Enum):
    lowest = 0
    highest = 1


class DisksRanking:
    def __init__(self, input_file_name):
        self.data_frame = pd.read_csv(input_file_name)

    def total_disks(self):
        pass

    def disks_per_datacenter(self):
        pass

    def age_for(self, old_young=None):
        if old_young is None:
            raise ValueError('No input if get age for youngest or oldest')
        if old_young == Ages.oldest:
            pass
        elif old_young == Ages.youngest:
            pass

    def youngest_disk(self):
        return self.age_for(Ages.youngest)

    def oldest_disk(self):
        return self.age_for(Ages.oldest)

    def average_age(self, datacenter=None):
        pass

    def average_ages_in_datacenters(self, data_centers):
        for data_center in data_centers:
            yield self.average_age(data_center)

    def average_read_write_io(self):
        pass

    def rank_read_write_io(self, low_high=None, disks_number=None):
        if low_high is None:
            raise ValueError('No input if check lowest or highest disks statistics')
        if disks_number is None:
            raise ValueError('No input how many disks has to be in ranking')

    def disks_with_errors(self):
        pass



if __name__ == '__main__':
    if len(sys.argv) == 2:
        input_file_name = sys.argv[1]
        if os.path.exists(input_file_name):
            disks_ranking = DisksRanking()
        else:
            raise IOError('Given input file doesnt exist')
    else:
        raise ValueError('No input filename given')



