type ProjectPaginationProps = {
    currentPage: number;
    totalPages: number;
    totalItems: number;
    firstItem: number;
    lastItem: number;
    onPageChange: (
        page: number,
    ) => void;
};

export function ProjectPagination({
    currentPage,
    totalPages,
    totalItems,
    firstItem,
    lastItem,
    onPageChange,
}: ProjectPaginationProps) {
    const hasItems = totalItems > 0;

    return (
        <footer className="project-pagination">
            <p className="project-pagination__result-count">
                {hasItems
                    ? `${firstItem}〜${lastItem}件 / 全${totalItems}件`
                    : '0件'}
            </p>

            <nav
                className="project-pagination__controls"
                aria-label="プロジェクト一覧のページネーション"
            >
                <button
                    className="project-pagination__button"
                    type="button"
                    disabled={currentPage <= 1}
                    aria-label="前のページ"
                    onClick={() => {
                        onPageChange(
                            currentPage - 1,
                        );
                    }}
                >
                    ‹
                </button>

                <span
                    className="project-pagination__current"
                    aria-current="page"
                >
                    {currentPage}
                </span>

                <button
                    className="project-pagination__button"
                    type="button"
                    disabled={
                        currentPage >= totalPages
                    }
                    aria-label="次のページ"
                    onClick={() => {
                        onPageChange(
                            currentPage + 1,
                        );
                    }}
                >
                    ›
                </button>
            </nav>
        </footer>
    );
}