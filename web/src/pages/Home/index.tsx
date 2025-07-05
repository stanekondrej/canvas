import { useLocation } from 'preact-iso';
import './style.css';

export function Home() {
    const { route } = useLocation()

    addEventListener("keydown", (e: KeyboardEvent) => {
        e.preventDefault();

        if (e.key !== ' ') {
            return;
        }

        route("/canvas");
    })

    return (
        <div>
            <h1>A global <span class="title">Canvas</span> for everyone</h1>

            <p>
                If you've ever wanted to draw on a very big sheet of paper with your
                friends, then try Canvas. Updates in real-time.
            </p>

            <a href="/canvas">Start drawing</a>
            <p><span class="hint">(Or, just press space)</span></p>
        </div>
    );
}

