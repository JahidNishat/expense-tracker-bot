export interface Transaction {
  id: number
  userId: number
  amount: number
  subcategoryId: string
  type: 'Expense' | 'Income' | 'Transfer'
  srcId: string
  dstId: string
  contactName: string
  timestamp: number
  remarks: string
  deletedAt: number
  createdAt: number
}

export interface Wallet {
  id: number
  userId: number
  type: 'Cash' | 'Bank'
  shortName: string
  name: string
  balance: number
  version: number
}

export interface Contact {
  id: number
  userId: number
  nickName: string
  fullName: string
  email: string
  netBalance: number
  lastTxnTimestamp: number
}

export interface BudgetStatus {
  id: number
  categoryId: string
  categoryName: string
  amount: number
  spent: number
  remaining: number
  percent: number
  alertAt: number
}

export interface BudgetAlert {
  categoryId: string
  categoryName: string
  budgetAmount: number
  spent: number
  percent: number
}

export interface CategorySpend {
  categoryId: string
  categoryName: string
  amount: number
  percent: number
}

export interface MonthlyComparison {
  month: string
  income: number
  expense: number
}

export interface MonthlyOverview {
  totalBalance: number
  monthIncome: number
  monthExpense: number
  budgetUsage: number
}

export interface TxnCategory {
  id: string
  name: string
}

export interface TxnSubcategory {
  id: string
  catId: string
  name: string
}

export interface Profile {
  id: number
  telegramId: number
  username: string
  firstName: string
  lastName: string
  timezone: string
  mobileNumber: string
}

export interface ChartData {
  overview: MonthlyOverview
  categories: CategorySpend[]
  comparison: MonthlyComparison[]
}
