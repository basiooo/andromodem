import type { JSX } from 'react'
import { useState } from 'react'

type Tab = {
  label: string
  key: string
  icon: JSX.Element
  content: JSX.Element
}

const useTabManagement = (initialTab: string) => {
  const [activeTab, setActiveTab] = useState(initialTab)

  const handleTabChange = (tabKey: string) => {
    setActiveTab(tabKey)
  }

  return {
    activeTab,
    handleTabChange
  }
}

export type { Tab }
export default useTabManagement