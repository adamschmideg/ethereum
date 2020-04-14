import pandas as pd
import matplotlib.pyplot as plt
from matplotlib.pyplot import figure

#url = "./"
url = "./data/"
converters = {'StartedAt': pd.to_datetime, 'FinishedAt': pd.to_datetime}
builds = pd.read_csv(url + "builds.csv", converters=converters)
jobs = pd.read_csv(url + "jobs.csv", converters=converters)
failures = pd.read_csv(url + "failures.csv")

bj = pd.merge(builds, jobs, left_on="Id", right_on="BuildId", suffixes=("_Build", "_Job"))
full = pd.merge(bj, failures, left_on="Id_Job", right_on="JobId")

cutoff_date = '2019-03-01'

df = full[['Test', 'PullRequestNumber']].groupby(by='Test').nunique()
not_one_pr = df[df.PullRequestNumber.ne(1)].Test.keys().values

df = full[full.StartedAt_Build >= cutoff_date] \
    [full.State_Job != "passed"] \
    [full.Test.isin(not_one_pr)]

#df = full[full.Test.isin(not_one_pr) & full.State_Job == "failed" & full.StartedAt_Build >= cutoff_date]
flaky = df[['Package', 'Test', 'Id_Job']] \
    .groupby(['Test']) \
    .agg({'Id_Job': ['first', 'count']})
flaky.columns = ['JobId', 'Count']
flaky.reset_index()
flaky['url'] = flaky.apply(lambda row: 'https://travis-ci.org/github/ethereum/go-ethereum/jobs/{}'.format(row.JobId), axis=1)
flaky = flaky.sort_values(['Count'], ascending=False)
pd.options.display.max_colwidth = 250


if 0:
    filtered = with_tests[with_tests.State_Job.eq('failed') & with_tests.StartedAt_Build.gt(cutoff_date)]
    cols = filtered[['Package', 'Test', 'Id_Job']]
    flaky = cols.groupby(['Test']).agg({'Id_Job': ['first', 'count']})
    flaky.columns = ['JobId', 'Count']
    flaky = flaky.reset_index()
    #flaky['url'] = df.apply(lambda row: 'https://travis-ci.org/github/ethereum/go-ethereum/jobs/{}'.format(row.JobId))

    failures_by_date = full[full.State_Job.eq('failed') & full.EventType.eq('push') & full.State_Job.gt(cutoff_date)][['Id_Build']].groupby(by=full.StartedAt_Build.dt.week).count().sort_values(by='StartedAt_Build', ascending=True)

    fig = plt.gcf()
    fig.set_size_inches(18.5, 10.5, forward=True)
    failures_by_date.plot()
    plt.show()