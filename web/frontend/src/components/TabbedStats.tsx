import { dailyStats, weeklyStats, monthlyStats } from "@/data/stats";
import { Download } from "lucide-react";
import { StatsChart } from "./StatsChart";
import { Button } from "./ui/button";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "./ui/tabs";
import { useEffect, useState } from "react";


const TabbedStats = () => {
  const [dailyStats, setDailyStats] = useState([])
  const [weeklyStats, setWeeklyStats] = useState([])
  const [monthlyStats, setMonthlyStats] = useState([])

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [dailyRes, weeklyRes, monthlyRes] = await Promise.all([
          fetch("http://localhost:8080/api/stats/daily"),
          fetch("http://localhost:8080/api/stats/weekly"),
          fetch("http://localhost:8080/api/stats/monthly"),
        ]);

        const [dailyData, weeklyData, monthlyData] = await Promise.all([
          dailyRes.json(),
          weeklyRes.json(),
          monthlyRes.json(),
        ]);

        setDailyStats(dailyData);
        setWeeklyStats(weeklyData);
        setMonthlyStats(monthlyData);
      } catch (err) {
        console.error("shit broke fetching stats", err);
      }
    };

    fetchStats();
  }, []);

  return (
    <Tabs defaultValue="daily" className="space-y-4">
      <div className="flex items-center justify-between">
        <TabsList className="">
          <TabsTrigger value="daily">Daily</TabsTrigger>
          <TabsTrigger value="weekly">Weekly</TabsTrigger>
          <TabsTrigger value="monthly">Monthly</TabsTrigger>
        </TabsList>
        <Button variant="outline" size="sm" className="gap-1">
          <Download className="h-4 w-4" />
          Export CSV
        </Button>
      </div>
      <TabsContent value="daily" className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <StatsChart
            title="Daily Hours"
            data={dailyStats}
            description="Hours tracked per day this week"
            dataKey="hours"
          />
          <StatsChart
            title="Daily Sessions"
            data={dailyStats}
            description="Sessions logged per day this week"
            dataKey="sessions"
          />
        </div>
      </TabsContent>
      <TabsContent value="weekly" className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <StatsChart
            title="Weekly Hours"
            data={weeklyStats}
            description="Hours tracked per week"
            dataKey="hours"
          />
          <StatsChart
            title="Weekly Sessions"
            data={weeklyStats}
            description="Sessions logged per week"
            dataKey="sessions"
          />
        </div>
      </TabsContent>
      <TabsContent value="monthly" className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <StatsChart
            title="Monthly Hours"
            data={monthlyStats}
            description="Hours tracked per month"
            dataKey="hours"
          />
          <StatsChart
            title="Monthly Sessions"
            data={monthlyStats}
            description="Sessions logged per month"
            dataKey="sessions"
          />
        </div>
      </TabsContent>
    </Tabs>
  )
}

export default TabbedStats

