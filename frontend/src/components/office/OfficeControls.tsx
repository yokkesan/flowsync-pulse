export function OfficeControls() {
    return (
        <footer className="office-controls">
            <button
                type="button"
                className="office-controls__button"
                disabled
            >
                マイク
            </button>

            <button
                type="button"
                className="office-controls__button"
                disabled
            >
                ステータス
            </button>
        </footer>
    );
}