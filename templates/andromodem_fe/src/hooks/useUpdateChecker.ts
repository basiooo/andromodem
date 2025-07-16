import { useEffect, useState } from "react"
import { toast } from "react-toastify"
import useSWR from "swr"

import { updateApi } from "@/api/updateApi"
import { config } from "@/config"
import type { GitHubRelease, UpdateInfo } from "@/types/update"

const useUpdateChecker = () => {
  const [updateInfo, setUpdateInfo] = useState<UpdateInfo | null>(null)
  const [isCheckingUpdate] = useState(false)

  const {
    data: latestRelease,
    error,
    isLoading
  } = useSWR<GitHubRelease>(
    "github_latest_release",
    updateApi.getLatestRelease,
    {
      revalidateIfStale: false,
      revalidateOnFocus: false,
      revalidateOnReconnect: false
    }
  )

  const compareVersions = (current: string, latest: string): boolean => {
    const cleanCurrent = current.replace(/^v/, "")
    const cleanLatest = latest.replace(/^v/, "")
    
    // Add debug logging
    console.log('Version comparison:', { current: cleanCurrent, latest: cleanLatest })
    
    if (cleanCurrent === "unknown") {
      console.log('Current version is unknown, showing update')
      return true
    }
    
    // Handle exact match first
    if (cleanCurrent === cleanLatest) {
      console.log('Versions are identical, no update needed')
      return false
    }
    
    const currentParts = cleanCurrent.split(".").map(Number)
    const latestParts = cleanLatest.split(".").map(Number)
    
    const maxLength = Math.max(currentParts.length, latestParts.length)
    
    for (let i = 0; i < maxLength; i++) {
      const currentPart = currentParts[i] || 0
      const latestPart = latestParts[i] || 0
      
      if (latestPart > currentPart) {
        return true
      }
      if (latestPart < currentPart) {
        return false
      }
    }
    
    return false
  }

  useEffect(() => {
    if (latestRelease && !isLoading && !error) {
      const currentVersion = config.VERSION
      const latestVersion = latestRelease.tag_name
      
      const hasUpdate = compareVersions(currentVersion, latestVersion)
      setUpdateInfo({
        hasUpdate,
        currentVersion,
        latestVersion,
        releaseInfo: latestRelease
      })
      
      if (hasUpdate) {
        toast.info(
          `Update available! Version ${latestVersion} has been released.`
        )
      }
    }
    
    if (error) {
      console.error("Failed to check for updates:", error)
    }
  }, [latestRelease, isLoading, error])

  return {
    updateInfo,
    isCheckingUpdate,
    isLoading,
    error
  }
}

export default useUpdateChecker