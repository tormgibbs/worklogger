import { Bar, BarChart, XAxis, YAxis } from 'recharts'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from './ui/card'

interface StatsChartProps {
  title: string
  description: string
  dataKey: string
  xAxisKey?: string
  color?: string
  data: any[]
}

export const StatsChart = ({
  title,
  description,
  data,
  dataKey,
  xAxisKey,
  color = 'blue',
}: StatsChartProps) => {
  const getXAxisKey = () => {
    if (xAxisKey) return xAxisKey
    if (data.length > 0) {
      const firstItem = data[0]
      if ('date' in firstItem) return 'date'
      if ('week_start' in firstItem) return 'week_start'
      if ('month' in firstItem) return 'month'
    }
    return 'name' // fallback
  }

  const chartData = data.length === 0 ? [{ [getXAxisKey()]: '', [dataKey]: 0 }] : data

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent>
        <BarChart
          width={400}
          height={250}
          data={chartData}
          margin={{
            top: 20,
            right: data.length === 1 ? 200 : 30,
            left: 20,
            bottom: 5
          }}
          barCategoryGap={data.length === 1 ? "20%" : "10%"}
        >
          <XAxis
            dataKey={getXAxisKey()}
            type="category"
            interval={0}
          />
          <YAxis />
          <Bar
            dataKey={dataKey}
            fill={data.length === 0 ? 'transparent' : color}
            maxBarSize={data.length === 1 ? 40 : 60} />
        </BarChart>
      </CardContent>
    </Card>
  )
}


