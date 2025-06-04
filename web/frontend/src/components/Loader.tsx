// components/Loader.tsx
export const Loader = () => (
  <div className="text-center text-sm text-muted-foreground">Loading...</div>
)

// components/ErrorMessage.tsx
export const ErrorMessage = ({ message }: { message: string }) => (
  <div className="text-red-500 text-sm text-center">{message}</div>
)
