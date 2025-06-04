import { Link, useMatchRoute } from '@tanstack/react-router'
import { type LinkProps } from '@tanstack/react-router'
import clsx from 'clsx'

type NavLinkProps = LinkProps & {
  activeClassName?: string
  inactiveClassName?: string
  exact?: boolean
  className?: string
}

export function NavLink({
  to,
  activeClassName = '',
  inactiveClassName = '',
  exact = false,
  className = '',
  ...rest
}: NavLinkProps) {
  const matchRoute = useMatchRoute()
  const isActive = !!matchRoute({ to, fuzzy: !exact })

  return (
    <Link
      to={to}
      {...rest}
      className={clsx(className, isActive ? activeClassName : inactiveClassName)}
      // className={[
      //   className,
      //   isActive ? activeClassName : inactiveClassName,
      // ].join(' ')}
    />
  )
}
