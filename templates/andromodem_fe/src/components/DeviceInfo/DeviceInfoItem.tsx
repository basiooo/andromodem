import type { ReactNode } from "react"

type Props = {
  icon: ReactNode;
  label: string;
  value: string | number;
};

const DeviceInfoItem = ({ icon, label, value }: Props) => (
  <li className="list-row">
    <div>{icon}</div>
    <div>
      <div>{label}</div>
      <div className="text-xs font-semibold opacity-60">{value || "-"}</div>
    </div>
  </li>
)

export default DeviceInfoItem
