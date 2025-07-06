import { ReactNode } from "preact/compat";

export const Link = (props: {
  href: string;
  target?: "_blank" | "_self";
  children: ReactNode;
}) => {
  return (
    <a href={props.href} target={props.target ?? "_self"} class="text-blue-700">
      {props.children}
    </a>
  );
};
