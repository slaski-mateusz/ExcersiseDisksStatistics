import json
from main import DisksRanking


if __name__ == '__main__':
    with open('test_data_1.json') as results_file:
        results = json.load(results_file)
    disks_ranking = DisksRanking('test_data_1.csv')
    for funct_name, expected_result in results.items():
        func_to_call = getattr(
            DisksRanking,
            funct_name
        )
        call_result = json.loads(
            func_to_call(
                disks_ranking
            )
        )
        if call_result == expected_result:
            print('%s passed' % funct_name)
        else:
            print('%s failed' % funct_name)
