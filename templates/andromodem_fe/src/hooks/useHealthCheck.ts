import useSWR from 'swr'

import { HealthApi } from '@/api/healthApi'
import type { BaseResponse } from '@/types/response'

const useHealthCheck = () => {
  const { data, error, isLoading, isValidating, mutate } = useSWR<BaseResponse>(
    'health_check',
    () => {
      return HealthApi.ping()
    },
    {
      revalidateIfStale: false,
      revalidateOnFocus: false,
      revalidateOnReconnect: false
    }
  )

  const healthError = error?.response?.status === 503 
    ? new Error(error.response?.data?.message || error.message || "Health check failed")
    : null

  return {
    data,
    error: healthError,
    isLoading,
    isValidating,
    mutate,
    isHealthy: data?.success === true
  }
}

export default useHealthCheck