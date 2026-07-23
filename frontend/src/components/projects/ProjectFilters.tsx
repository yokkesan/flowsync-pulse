import type {
    ProjectStatusFilter,
} from '../../types/project';

type ProjectFiltersProps = {
    searchKeyword: string;
    statusFilter: ProjectStatusFilter;
    onSearchKeywordChange: (
        value: string,
    ) => void;
    onStatusFilterChange: (
        value: ProjectStatusFilter,
    ) => void;
};

export function ProjectFilters({
    searchKeyword,
    statusFilter,
    onSearchKeywordChange,
    onStatusFilterChange,
}: ProjectFiltersProps) {
    return (
        <div className="project-filters">
            <label className="project-filters__search">
                <span
                    className="project-filters__search-icon"
                    aria-hidden="true"
                >
                    ⌕
                </span>

                <input
                    className="project-filters__search-input"
                    type="search"
                    value={searchKeyword}
                    placeholder="プロジェクト名で検索"
                    onChange={(event) => {
                        onSearchKeywordChange(
                            event.target.value,
                        );
                    }}
                />
            </label>

            <label className="project-filters__status">
                <span className="project-filters__status-label">
                    ステータス
                </span>

                <select
                    className="project-filters__status-select"
                    value={statusFilter}
                    onChange={(event) => {
                        onStatusFilterChange(
                            event.target
                                .value as ProjectStatusFilter,
                        );
                    }}
                >
                    <option value="all">
                        すべて
                    </option>

                    <option value="active">
                        進行中
                    </option>

                    <option value="completed">
                        完了
                    </option>

                    <option value="paused">
                        停止中
                    </option>
                </select>
            </label>
        </div>
    );
}