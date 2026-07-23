type TaskPaginationProps = {
    currentPage: number;
    totalPages: number;
    totalItems: number;
    itemsPerPage: number;
    onPageChange: (
        page: number,
    ) => void;
};

export function TaskPagination({
    currentPage,
    totalPages,
    totalItems,
    itemsPerPage,
    onPageChange,
}: TaskPaginationProps) {
    const startItem =
        totalItems === 0
            ? 0
            : (currentPage - 1) *
                  itemsPerPage +
              1;

    const endItem = Math.min(
        currentPage * itemsPerPage,
        totalItems,
    );

    return (
        <footer className="task-pagination">
            <p className="task-pagination__summary">
                全{totalItems}件中{' '}
                {startItem}～{endItem}
                件を表示
            </p>

            <div className="task-pagination__controls">
                <button
                    className="task-pagination__button"
                    type="button"
                    disabled={currentPage <= 1}
                    onClick={() => {
                        onPageChange(
                            currentPage - 1,
                        );
                    }}
                >
                    前へ
                </button>

                <span className="task-pagination__current">
                    {currentPage} / {totalPages}
                </span>

                <button
                    className="task-pagination__button"
                    type="button"
                    disabled={
                        currentPage >= totalPages
                    }
                    onClick={() => {
                        onPageChange(
                            currentPage + 1,
                        );
                    }}
                >
                    次へ
                </button>
            </div>
        </footer>
    );
}