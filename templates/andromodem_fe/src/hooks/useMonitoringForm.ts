import { zodResolver } from "@hookform/resolvers/zod"
import { useEffect, useRef,useState } from "react"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { type MonitoringConfig,MonitoringMethod, type MonitoringMethodValue } from "@/types/monitoring"

const monitoringSchema = z.object({
  method: z.enum([MonitoringMethod.HTTP, MonitoringMethod.HTTPS, MonitoringMethod.WS, MonitoringMethod.PING, MonitoringMethod.PING_BY_DEVICE]),
  host: z.string()
    .min(1, "Host is required")
    .refine((val) => !val.startsWith('http://') && !val.startsWith('https://'), {
      message: "Host should not include protocol (http:// or https://)"
    }),
  max_failures: z.number()
    .min(1, "Max failures must be at least 1")
    .max(100, "Max failures cannot exceed 100"),
  checking_interval: z.number()
    .min(5, "Checking interval must be at least 5 seconds")
    .max(3600, "Checking interval cannot exceed 1 hour"),
  airplane_mode_delay: z.number()
    .min(1, "Airplane mode delay cannot 0")
    .max(3600, "Airplane mode delay cannot exceed 1 hour")
})

export type MonitoringFormData = z.infer<typeof monitoringSchema>;

interface UseMonitoringFormProps {
  monitoringConfig: MonitoringConfig | null;
  onSubmit: (data: MonitoringFormData) => Promise<void>;
  onCancel: () => void;
}

export const useMonitoringForm = ({
  monitoringConfig,
  onSubmit,
  onCancel
}: UseMonitoringFormProps) => {
  const [currentMethod, setCurrentMethod] = useState<MonitoringMethodValue>(MonitoringMethod.HTTP)
  const originalDataRef = useRef<MonitoringFormData | null>(null)

  const {
    register,
    handleSubmit,
    setValue,
    reset,
    getValues,
    formState: { errors, isSubmitting }
  } = useForm<MonitoringFormData>({
    resolver: zodResolver(monitoringSchema),
    defaultValues: {
      method: MonitoringMethod.HTTP,
      host: '',
      max_failures: 3,
      checking_interval: 30,
      airplane_mode_delay: 5,
    },
    mode: 'onChange'
  })

  useEffect(() => {
    if (monitoringConfig) {
      const originalData: MonitoringFormData = {
        method: monitoringConfig.method,
        host: monitoringConfig.host || '',
        max_failures: monitoringConfig.max_failures || 3,
        checking_interval: monitoringConfig.checking_interval || 60,
        airplane_mode_delay: monitoringConfig.airplane_mode_delay || 0
      }
      
      originalDataRef.current = originalData
      
      setValue('method', originalData.method)
      setValue('host', originalData.host)
      setValue('max_failures', originalData.max_failures)
      setValue('checking_interval', originalData.checking_interval)
      setValue('airplane_mode_delay', originalData.airplane_mode_delay)
      setCurrentMethod(originalData.method)

    } else {
      const defaultData: MonitoringFormData = {
        method: MonitoringMethod.HTTP,
        host: '',
        max_failures: 3,
        checking_interval: 30,
        airplane_mode_delay: 5,
      }
      
      originalDataRef.current = defaultData
      
      setValue('method', defaultData.method)
      setValue('host', defaultData.host)
      setValue('max_failures', defaultData.max_failures)
      setValue('checking_interval', defaultData.checking_interval)
      setValue('airplane_mode_delay', defaultData.airplane_mode_delay)
      setCurrentMethod(defaultData.method)
    }
  }, [monitoringConfig, setValue, getValues])

  const handleFormSubmit = handleSubmit(async (data: MonitoringFormData) => {
    await onSubmit(data)
  })

  const handleCancel = () => {
    
    if (originalDataRef.current) {
      const originalData = originalDataRef.current      
      setValue('method', originalData.method)
      setValue('host', originalData.host)
      setValue('max_failures', originalData.max_failures)
      setValue('checking_interval', originalData.checking_interval)
      setValue('airplane_mode_delay', originalData.airplane_mode_delay)
      setCurrentMethod(originalData.method)
    }
    
    onCancel()
  }

  return {
    register,
    handleSubmit: handleFormSubmit,
    setValue,
    reset,
    errors,
    isSubmitting,
    currentMethod,
    setCurrentMethod,
    handleCancel
  }
}