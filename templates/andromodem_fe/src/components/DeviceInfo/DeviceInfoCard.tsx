import type { CSSProperties,ReactNode } from "react"

import { showModal } from "@/utils/common"

type Props = {
  title: string;
  value: number;
  modal_id: string;
  reverse_radial_color?: boolean;
  children: ReactNode;
};

const getColorClass = (value: number, reverse: boolean): string => {
  const thresholds = reverse
    ? [
        { threshold: 80, color: "text-red-400" },
        { threshold: 50, color: "text-orange-400" },
        { threshold: 15, color: "text-blue-400" },
        { threshold: 0, color: "text-green-400" }
      ]
    : [
        { threshold: 80, color: "text-green-400" },
        { threshold: 50, color: "text-blue-400" },
        { threshold: 15, color: "text-orange-400" },
        { threshold: 0, color: "text-red-400" }
      ]

  return thresholds.find(t => value > t.threshold)?.color ?? "text-gray-400"
}

const DeviceInfoCard = ({
  title,
  value,
  modal_id,
  reverse_radial_color = false,
  children
}: Props) => {
  const colorClass = getColorClass(value, reverse_radial_color)

  return (
    <div className="card bg-base-100 w-full shadow-sm">
      <div className="card-body">
        <h2 className="card-title">{title}</h2>
        <div className="flex items-center justify-around gap-4">
          <div className="bg-base-300 rounded-full p-2">
            <div
              className={`radial-progress ${colorClass}`}
              style={
                {
                  "--value": value,
                  "--size": "4rem"
                } as CSSProperties
              }
              aria-valuenow={value}
              role="progressbar"
            >
              {value}%
            </div>
          </div>
          <div className="md:text-xs">{children}</div>
        </div>
        <div className="text-center">
          <button
            className="btn btn-sm btn-info"
            onClick={() => showModal(modal_id)}
          >
            Detail
          </button>
        </div>
      </div>
    </div>
  )
}

export default DeviceInfoCard
