# (c) Mateusz Slaski for OVHcloud

import sys
import os
import pandas as pd
import json

pd.options.display.max_columns = None
pd.options.display.max_colwidth = 0
pd.options.display.float_format = '{:.2f}'.format

input_columns = [
    'datacenter',
    'hostname',
    'disk_serial',
    'disk_age_sec',
    'total_reads',
    'total_writes',
    'ave_io_lat',
    'tot_uncr_r_er',
    'tot_uncr_w_er'
]


def sec_to_days(seconds):
    return seconds / 86400


class DisksRanking:
    def __init__(self, input_filename):
        self.disks_data = pd.read_csv(
            filepath_or_buffer=input_filename,
            delimiter=';',
            header=None
        )
        self.disks_data.columns = input_columns
        if self.check_for_duplicated_disks():
            raise ValueError('There are duplicated serial numbers of disks in data')
        self.disks_data['sum_read_write'] = self.disks_data['total_reads'] + self.disks_data['total_writes']
        self.disks_data['disk_age_days'] = pd.Series(
            sec_to_days(
                self.disks_data['disk_age_sec']
            )
        )

    def check_for_duplicated_disks(self):
        have_duplicated_serials = False
        for dupl_check in self.disks_data.disk_serial.duplicated():
            if dupl_check:
                have_duplicated_serials = True
        return have_duplicated_serials

    def total_disks(self):
        return json.dumps(
            {
                'total_disks': len(self.disks_data.index)
            }
        )

    def disks_per_datacenter(self):
        return json.dumps(
            self.disks_data.groupby('datacenter')['datacenter'].count().to_dict()
        )

    def youngest_disk(self):
        disk_data = self.disks_data[self.disks_data['disk_age_sec'] == self.disks_data['disk_age_sec'].min()][[
            'datacenter',
            'hostname',
            'disk_serial',
            'disk_age_days'
        ]].to_dict('records')
        return json.dumps(
            disk_data
        )

    def oldest_disk(self):
        disk_data = self.disks_data[self.disks_data['disk_age_sec'] == self.disks_data['disk_age_sec'].max()][[
            'datacenter',
            'hostname',
            'disk_serial',
            'disk_age_days'
        ]].to_dict('records')
        return json.dumps(
            disk_data
        )

    def average_age_per_dc(self):
        avgs_dc = self.disks_data.groupby('datacenter')['disk_age_days'].mean().to_dict()
        return json.dumps(avgs_dc)

    def average_read_write(self):
        return json.dumps(
            {
                'avg_reads': self.disks_data['total_reads'].mean(),
                'avg_writes': self.disks_data['total_writes'].mean()
            }
        )

    def rank_read_write_io(self, top=None, disks_number=None):
        if top is None:
            raise ValueError('No input if check lowest or highest disks statistics')
        if disks_number is None:
            raise ValueError('No input how many disks has to be in ranking')
        return json.dumps(
            self.disks_data.sort_values(
                'sum_read_write',
                ascending=not top
            ).head(disks_number)[[
                'datacenter',
                'hostname',
                'disk_serial',
                'total_reads',
                'total_writes',
                'sum_read_write'
            ]].to_dict('records')
        )

    def most_loaded_5_disks(self):
        return self.rank_read_write_io(top=True, disks_number=5)

    def less_loaded_5_disks(self):
        return self.rank_read_write_io(top=False, disks_number=5)

    def disks_with_errors(self):
        dwe = self.disks_data.loc[
            (self.disks_data['tot_uncr_r_er'] > 0) |
            (self.disks_data['tot_uncr_w_er'] > 0)
            ][[
                'datacenter',
                'hostname',
                'disk_serial',
                'tot_uncr_r_er',
                'tot_uncr_w_er'
            ]]
        return json.dumps(dwe.to_dict('records'))


if __name__ == '__main__':
    if len(sys.argv) == 2:
        input_file_name = sys.argv[1]
        if os.path.exists(input_file_name):
            disks_ranking = DisksRanking(input_file_name)
            print('\nTotal disks')
            print(disks_ranking.total_disks())
            print('\nDisks per datacenter')
            print(disks_ranking.disks_per_datacenter())
            print('\nYoungest disk')
            print(disks_ranking.youngest_disk())
            print('\noldest disk')
            print(disks_ranking.oldest_disk())
            print('\nAverage age per data center')
            print(disks_ranking.average_age_per_dc())
            print('\nAverage read write')
            print(disks_ranking.average_read_write())
            print('\n5 most loaded disks')
            print(disks_ranking.most_loaded_5_disks())
            print('\n5 less loaded disks')
            print(disks_ranking.less_loaded_5_disks())
            print('\nDisks with read or write errors')
            print(disks_ranking.disks_with_errors())
        else:
            raise IOError('Given input file doesnt exist')
    else:
        raise ValueError('No input filename given')
