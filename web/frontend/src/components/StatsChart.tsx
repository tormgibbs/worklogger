import { Bar, BarChart, Tooltip, XAxis, YAxis } from 'recharts'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from './ui/card'
import type { NameType, ValueType } from 'recharts/types/component/DefaultTooltipContent'

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
    return 'name'
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
            tickFormatter={(value, index) => formatXAxis(value, getXAxisKey(), index)}
          />
          <YAxis />
          <Tooltip
            formatter={(value, name) => formatTooltipValue(value, name)}
            labelFormatter={(label) => formatTooltipLabel(label, getXAxisKey())}
          />
          <Bar
            dataKey={dataKey}
            fill={data.length === 0 ? 'transparent' : color}
            maxBarSize={data.length === 1 ? 40 : 60} />
        </BarChart>
      </CardContent>
    </Card>
  )
}

const formatXAxis = (value: string, key: string, index: number) => {
  if (!value) return ''

  if (key === 'date') {
    return new Date(value).toLocaleDateString(undefined, { weekday: 'short' }) // "Mon", "Tue"
  }

  if (key === 'week_start') {
    return `Week ${index + 1}`
  }

  if (key === 'month') {
    const [year, month] = value.split('-')
    const date = new Date(Number(year), Number(month) - 1)
    return date.toLocaleDateString(undefined, { month: 'short' }) // "Jun"
  }

  return value
}

const formatTooltipLabel = (label: string, key: string) => {
  if (!label) return ''

  if (key === 'date') {
    return new Date(label).toLocaleDateString(undefined, {
      weekday: 'long',
      day: 'numeric',
      month: 'short',
    })
  }

  if (key === 'week_start') {
    return `Starting ${new Date(label).toLocaleDateString(undefined, {
      day: 'numeric',
      month: 'short',
    })}`
  }

  if (key === 'month') {
    const [year, month] = label.split('-')
    const date = new Date(Number(year), Number(month) - 1)
    return date.toLocaleDateString(undefined, {
      month: 'long',
      year: 'numeric',
    })
  }

  return label
}

const formatTooltipValue = (value: ValueType, name: NameType) => {
  const unit = name === 'hours' ? 'hrs' : 'sessions'
  const capitalized = name.charAt(0).toUpperCase() + name.slice(1)
  return [`${value} ${unit}`, capitalized]
}





