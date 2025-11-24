export const SIDEBAR_NAVIIGATION = [
  {
    title: 'Dashboard',
    to: { name: 'dashboard' },
    icon: 'lucide:layout-dashboard',
  },
  {
    title: 'Accounts',
    to: { name: 'dashboard' },
    icon: 'lucide:banknote',
  },
  {
    title: 'Categories',
    to: { name: 'dashboard' },
    icon: 'lucide:layout-grid',
  },
  {
    title: 'Transactions',
    to: { name: 'dashboard' },
    icon: 'lucide:arrow-up-down',
  },
  {
    title: 'Shared Expenses',
    to: { name: 'dashboard' },
    icon: 'lucide:users',
  },

  {
    title: 'Subscriptions',
    to: { name: 'dashboard' },
    icon: 'lucide:repeat',
  },
  {
    title: 'Goals',
    to: { name: 'dashboard' },
    icon: 'lucide:goal',
  },
  // {
  //   title: 'Budget',
  //   to: { name: 'budget' },
  //   icon: 'lucide:pie-chart',
  // },
  // {
  //   title: 'Loans',
  //   to: { name: 'loans' },
  //   icon: 'lucide:hand-coins',
  // },
  // {
  //   title: 'Templates',
  //   to: { name: 'templates' },
  //   icon: 'lucide:book-template',
  // },
]

export const SIDEBAR_FOOTER_NAVIGATION = [
  {
    title: 'Settings',
    to: { name: 'dashboard' },
    icon: 'lucide:settings',
  },
]

export const SETTINGS_NAVIGATION = [
  {
    title: 'Profile',
    icon: 'lucide:user',
    to: { name: 'settings-profile' },
    description: 'Personal information',
  },
  {
    title: 'Security',
    icon: 'lucide:shield',
    to: { name: 'settings-security' },
    description: 'Password and authentication',
  },
  {
    title: 'Connections',
    icon: 'lucide:link',
    to: { name: 'settings-connections' },
    description: 'Email addresses and social accounts',
  },
  // {
  //   title: 'Notifications',
  //   icon: 'lucide:bell',
  //   to: { name: 'settings-notifications' },
  //   description: 'Email and alerts',
  // },
  // {
  //   title: 'Data and privacy',
  //   icon: 'lucide:lock',
  //   to: { name: 'settings-privacy' },
  //   description: 'Data and privacy',
  // },
]
