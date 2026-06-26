import { redirect } from "next/navigation";

// Registration on Vouch is GitHub-OAuth only — there is no separate sign-up.
export default function RegisterPage() {
  redirect("/login");
}
