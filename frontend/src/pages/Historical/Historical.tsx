import { HStack, Tab, TabList, TabPanel, TabPanels, Tabs, Text, VStack } from '@chakra-ui/react'
import { utils } from 'ethers'
import { useEffect, useState } from 'react'
import { Tooltips } from '../../config'
import { useEthereum } from '../../contexts/EthereumContext'
import { layoutConfig } from '../../layoutConfig'
import { Percentiles, TotalsWithId } from '../../libs/ethereum'
import { HistoricalChart } from './HistoricalChart'
import { ChartData, ChartDataBucket, TimeBucket } from './HistoricalTypes'

interface ChartRange {
  hour: ChartDataBucket
  day: ChartDataBucket,
  month: ChartDataBucket
}


function TotalTabPanel({ bucket }: { bucket: ChartDataBucket | undefined }) {
  return (
    <VStack spacing={0} gridGap={layoutConfig.gap} mt={8} align="flex-start">
      <HistoricalChart
        title="BaseFee"
        dataKey={["baseFee"]}
        percentilesKey="baseFeePercentiles"
        tooltip={Tooltips.baseFeeMedian}
        bucket={bucket} />
      <HistoricalChart
        title="Burned"
        dataKey={["burned"]}
        tooltip={Tooltips.burned}
        bucket={bucket} />
      <HistoricalChart
        title="Net Issuance"
        dataKey={["issuance"]}
        tooltip={Tooltips.netIssuance}
        bucket={bucket} />
      <HistoricalChart
        title="Rewards and Tips"
        dataKey={["rewards", "tips"]}
        tooltip={`${Tooltips.rewards} ${Tooltips.tips}`}
        bucket={bucket} />
    </VStack>
  )
}

function TabTitle({ bucket }: { bucket: ChartDataBucket | undefined }) {
  return (
    <Text flex={1} pr={4} align="right" variant="brandSecondary">
      {bucket && (<>{bucket.data.length} {bucket.type}s</>)}
      {!bucket && (<>Rendering... please wait!</>)}
    </Text>
  )
}

export function Historical() {
  const ethereum = useEthereum()
  const [data, setData] = useState<ChartRange>()
  const [bucket, setBucket] = useState<ChartDataBucket>()

  useEffect(() => {
    if (!ethereum.eth) {
      return
    }

    const init = async () => {
      const response = await ethereum.eth!.getInitialAggregatesData()
      const formatTimestampToDateString = (id: string, bucket: TimeBucket) => {
        const startTime = id.split(':')
        const date = new Date(parseInt(startTime[0]) * 1000)

        switch (bucket) {
          case 'hour':
            return new Intl.DateTimeFormat(navigator.language, {
              hour: 'numeric',
              day: 'numeric',
              month: 'short',
              timeZone: 'UTC'
            }).format(date)
          case 'day':
            return new Intl.DateTimeFormat(navigator.language, {
              day: 'numeric',
              month: 'short',
              timeZone: 'UTC'
            }).format(date)
          case 'month':
            return new Intl.DateTimeFormat(navigator.language, {
              month: 'long',
              timeZone: 'UTC'
            }).format(date)
        }
      }

      // Since the chart software adds all the stacked bars, it really doesn't make
      // sense since we are not getting the total of base fee, so for the ones we are
      // rendering, just subract it so the number shows the maximum base fee.
      const normalizePercentiles = (percentiles: Percentiles): Percentiles => {
        return {
          Minimum: percentiles.Minimum,
          Median: percentiles.Median - percentiles.Minimum,
          Maximum: percentiles.Maximum,
          ninetieth: percentiles.ninetieth - percentiles.Median,
        }
      }

      const formatToChartData = (totals: TotalsWithId[], bucket: TimeBucket) => totals.map<ChartData>(total => (
        {
          timestamp: formatTimestampToDateString(total.id, bucket),
          baseFee: total.baseFee,
          baseFeePercentiles: normalizePercentiles(total.baseFeePercentiles),
          burned: parseFloat(utils.formatUnits(total.burned, 'ether')),
          issuance: parseFloat(utils.formatUnits(total.issuance, 'ether')),
          rewards: parseFloat(utils.formatUnits(total.rewards, 'ether')),
          tips: parseFloat(utils.formatUnits(total.tips, 'ether')),
          netReduction: total.netReduction,
        }
      ))

      const mutatedData: ChartRange = {
        hour: { type: "hour", data: formatToChartData(response.totalsPerHour, 'hour').reverse() },
        day: { type: "day", data: formatToChartData(response.totalsPerDay, 'day').reverse() },
        month: { type: "month", data: formatToChartData(response.totalsPerMonth, 'month').reverse() },
      }

      setData(mutatedData);
      setBucket(mutatedData.hour)
    }

    init()
  }, [ethereum])

  return (
    <Tabs variant="inline" isLazy onChange={(index) => data && setBucket(index === 0 ? data.hour : index === 1 ? data.day : data.month)}>
      <HStack>
        <TabTitle bucket={bucket} />
        <TabList display="inline-flex">
          <Tab>Hour</Tab>
          <Tab>Day</Tab>
          <Tab>Month</Tab>
        </TabList>
      </HStack>
      <TabPanels>
        <TabPanel><TotalTabPanel bucket={data && data.hour} /></TabPanel>
        <TabPanel><TotalTabPanel bucket={data && data.day} /></TabPanel>
        <TabPanel><TotalTabPanel bucket={data && data.month} /></TabPanel>
      </TabPanels>
    </Tabs>
  )
}
