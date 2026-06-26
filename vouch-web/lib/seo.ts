import type { Metadata } from "next";

const BASE_URL = process.env.NEXT_PUBLIC_BASE_URL ?? "https://vouch.dev";
const SITE_NAME = "Vouch";
const DEFAULT_DESCRIPTION =
  "Demand-first builder reputation. Companies post problems, builders ship solutions, the community vouches.";

interface PageSEOOptions {
  title: string;
  description?: string;
  path?: string;
  ogImage?: string;
  noIndex?: boolean;
}

export function buildMetadata({
  title,
  description = DEFAULT_DESCRIPTION,
  path = "",
  ogImage,
  noIndex = false,
}: PageSEOOptions): Metadata {
  const url = `${BASE_URL}${path}`;
  const image = ogImage ?? `${BASE_URL}/og.png`;

  return {
    title: `${title} | ${SITE_NAME}`,
    description,
    metadataBase: new URL(BASE_URL),
    alternates: { canonical: url },
    robots: noIndex ? { index: false, follow: false } : { index: true, follow: true },
    openGraph: {
      title,
      description,
      url,
      siteName: SITE_NAME,
      images: [{ url: image, width: 1200, height: 630, alt: title }],
      type: "website",
    },
    twitter: {
      card: "summary_large_image",
      title,
      description,
      images: [image],
    },
  };
}

export function builderProfileMetadata(username: string, bio?: string): Metadata {
  return buildMetadata({
    title: `@${username}`,
    description: bio ?? `View ${username}'s builder profile on Vouch.`,
    path: `/u/${username}`,
  });
}
