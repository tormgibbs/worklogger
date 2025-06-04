import { LayoutDashboard, ListTodo } from 'lucide-react'
import { NavLink } from './NavLink' // adjust the path if needed

export default function Header() {
  return (
    <header className="p-4 flex gap-2 bg-white text-black items-center justify-between">
      <p className="text-2xl font-bold">Worklogger</p>
      <nav className="flex flex-row gap-4">
        <NavLink
          to="/dashboard"
          className="flex items-center gap-2 font-bold"
          activeClassName="border p-2 rounded-lg text-blue-600"
          inactiveClassName="text-gray-600 hover:text-black"
        >
          <LayoutDashboard size={20} />
          Dashboard
        </NavLink>
        <NavLink
          to="/sessions"
          className="flex items-center gap-2 font-bold"
          activeClassName="border p-2 rounded-lg text-blue-600"
          inactiveClassName="text-gray-600 hover:text-black"
        >
          <ListTodo size={20} />
          Sessions
        </NavLink>
      </nav>
    </header>
  )
}
