import { Navbar } from "@/components/ui/navbar";

// Shell for the authenticated/app pages — landing ("/") opts out of this by
// living outside the (app) group, so it can render its own full-bleed chrome.
export default function AppLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      <Navbar />
      <main className="mx-auto max-w-6xl px-4 py-10">{children}</main>
    </>
  );
}
