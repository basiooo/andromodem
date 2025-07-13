type Props = {
  count?: number;
};

const DeviceInfoItemSkeleton = ({ count = 5 }: Props) => {
  return (
    <ul className="list bg-base-200 rounded-box shadow-md my-5">
      <li className="p-4 pb-2 text-xl tracking-wide">Device</li>
      {Array.from({ length: count }).map((_, index) => (
        <li className="list-row animate-pulse" key={index}>
          <div className="w-6 h-6 bg-base-300 rounded-full" />
          <div className="flex flex-col gap-1">
            <div className="w-28 h-3 bg-base-300 rounded" />
            <div className="w-24 h-2 bg-base-300 rounded opacity-60" />
          </div>
        </li>
      ))}
    </ul>
  )
}

export default DeviceInfoItemSkeleton
