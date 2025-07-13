import type { StorageUnit } from "@/types/utils"

export const getPercentage = (
    totalValue: number,
    currentValue: number
): number => {
    return Math.floor((currentValue / totalValue) * 100)
}

export const UnitFactors: Record<StorageUnit, number> = {
  KB: 1,
  MB: 1024,
  GB: 1024 * 1024
}
export const convertStorageUnit = (
  value: number,
  from: StorageUnit,
  to: StorageUnit,
  precision: number = 2,
  roundDown: boolean = false
): number => {
  const bytes = value * (UnitFactors[from] || 1)
  const result = bytes / (UnitFactors[to] || 1)
  const rounded = roundDown ? Math.floor(result) : parseFloat(result.toFixed(precision))
  return rounded
}

export const secondsForHuman = (totalSeconds: number) => {
    const days = Math.floor(totalSeconds / 86400)
    const hours = Math.floor((totalSeconds % 86400) / 3600)
    const minutes = Math.floor((totalSeconds % 3600) / 60)
    const secs = totalSeconds % 60

    return `${days > 0 ? `${days} day${days !== 1 ? "s" : ""} ` : ""}${
        hours > 0 || days > 0 ? `${hours} hour${hours !== 1 ? "s" : ""} ` : ""
    }${minutes > 0 || hours > 0 || days > 0 ? `${minutes} minute${minutes !== 1 ? "s" : ""} ` : ""}${secs} second${secs !== 1 ? "s" : ""}`
}