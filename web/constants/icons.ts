export interface IconCategory {
  category: string
  icons: string[]
}

export const AVAILABLE_ICONS: IconCategory[] = [
  {
    category: 'Finance',
    icons: [
      'lucide:wallet',
      'lucide:banknote',
      'lucide:landmark',
      'lucide:credit-card',
      'lucide:piggy-bank',
      'lucide:coins',
      'lucide:building-2',
      'lucide:safe',
      'lucide:receipt',
      'lucide:trending-up',
      'lucide:trending-down',
      'lucide:dollar-sign',
      'lucide:currency',
    ],
  },
  {
    category: 'Business',
    icons: [
      'lucide:briefcase',
      'lucide:building',
      'lucide:building-2',
      'lucide:briefcase-business',
      'lucide:handshake',
      'lucide:target',
      'lucide:chart-bar',
      'lucide:pie-chart',
    ],
  },
  {
    category: 'Shopping',
    icons: [
      'lucide:shopping-bag',
      'lucide:shopping-cart',
      'lucide:store',
      'lucide:package',
      'lucide:gift',
    ],
  },
  {
    category: 'Home',
    icons: [
      'lucide:home',
      'lucide:house',
      'lucide:key',
      'lucide:door-open',
      'lucide:sofa',
      'lucide:lightbulb',
      'lucide:wrench',
    ],
  },
  {
    category: 'Health',
    icons: [
      'lucide:heart',
      'lucide:activity',
      'lucide:stethoscope',
      'lucide:pill',
      'lucide:cross',
      'lucide:bandage',
      'lucide:syringe',
    ],
  },
  {
    category: 'Transport',
    icons: [
      'lucide:car',
      'lucide:train',
      'lucide:plane',
      'lucide:ship',
      'lucide:bike',
      'lucide:fuel',
      'lucide:map-pin',
    ],
  },
  {
    category: 'Food',
    icons: [
      'lucide:utensils',
      'lucide:utensils-crossed',
      'lucide:coffee',
      'lucide:cup-soda',
      'lucide:cookie',
      'lucide:cherry',
    ],
  },
  {
    category: 'Entertainment',
    icons: [
      'lucide:film',
      'lucide:tv',
      'lucide:gamepad-2',
      'lucide:music',
      'lucide:headphones',
      'lucide:theater',
      'lucide:palette',
    ],
  },
  {
    category: 'Education',
    icons: [
      'lucide:graduation-cap',
      'lucide:book',
      'lucide:book-open',
      'lucide:school',
      'lucide:laptop',
      'lucide:pen-tool',
    ],
  },
  {
    category: 'Utilities',
    icons: [
      'lucide:zap',
      'lucide:wifi',
      'lucide:phone',
      'lucide:smartphone',
      'lucide:cloud',
      'lucide:battery',
    ],
  },
  {
    category: 'Travel',
    icons: [
      'lucide:map',
      'lucide:compass',
      'lucide:camera',
      'lucide:suitcase',
      'lucide:hotel',
      'lucide:plane-takeoff',
    ],
  },
  {
    category: 'Other',
    icons: [
      'lucide:star',
      'lucide:heart-handshake',
      'lucide:users',
      'lucide:user',
      'lucide:calendar',
      'lucide:clock',
      'lucide:bell',
      'lucide:settings',
    ],
  },
]
