import { ResponsiveContainer } from "recharts";

// https://github.com/recharts/recharts/issues/1767#issuecomment-598607012
export function CustomResponsiveContainer(props: any) {
  return (
    <div style={{ width: '100%', height: '100%', position: 'relative' }}>
      <div
        style={{
          width: '100%',
          height: '100%',
          position: 'absolute',
          top: 0,
          left: 0
        }}
      >
        <ResponsiveContainer {...props} />
      </div>
    </div>
  );
}