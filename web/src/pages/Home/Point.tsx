import { ReactNode } from "preact/compat";

export const Point = (props: {
  image: string;
  title: string;
  children: ReactNode;
}) => {
  return (
    <div class="flex flex-col items-center gap-8 px-4">
      <img
        src={props.image}
        class="w-32 drop-shadow-xl hover:drop-shadow-blue-500 hover:-translate-y-1 transition-all"
      />

      <h2 class="text-2xl text-black font-bold">{props.title}</h2>

      <p class="text-justify">{props.children}</p>
    </div>
  );
};
