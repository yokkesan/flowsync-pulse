import type {
    WorkContextChangedPayload,
} from '../../types/realtime';

type OfficeUserDetailPanelProps = {
    displayName: string;
    workContext: WorkContextChangedPayload | null;
};

export function OfficeUserDetailPanel({
    displayName,
    workContext,
}: OfficeUserDetailPanelProps) {
    return (
        <aside className="office-user-detail-panel">
            <div className="office-user-detail-panel__header">
                <p className="office-user-detail-panel__eyebrow">
                    選択中のユーザー
                </p>

                <h2 className="office-user-detail-panel__name">
                    {displayName}
                </h2>
            </div>

            <div className="office-user-detail-panel__content">
                {workContext ? (
                    <dl className="office-user-detail-panel__list">
                        <div className="office-user-detail-panel__item">
                            <dt>プロジェクト</dt>
                            <dd>
                                {workContext.project_name ??
                                    '未設定'}
                            </dd>
                        </div>

                        <div className="office-user-detail-panel__item">
                            <dt>タスク</dt>
                            <dd>
                                {workContext.task_name ??
                                    workContext.task_key ??
                                    workContext.ticket_key ??
                                    '未設定'}
                            </dd>
                        </div>

                        <div className="office-user-detail-panel__item">
                            <dt>ブランチ</dt>
                            <dd>
                                {workContext.branch_name ??
                                    '未設定'}
                            </dd>
                        </div>

                        <div className="office-user-detail-panel__item">
                            <dt>リポジトリ</dt>
                            <dd>
                                {workContext.repository_name ??
                                    '未設定'}
                            </dd>
                        </div>
                    </dl>
                ) : (
                    <p className="office-user-detail-panel__empty">
                        作業情報はありません。
                    </p>
                )}
            </div>
        </aside>
    );
}