import {
    useEffect,
    useMemo,
    useState,
    type FormEvent,
} from 'react';

import type {
    ProjectStatus,
    ProjectWriteRequest,
} from '../../types/project';
import type { CompanyUser } from '../../types/user';

export type ProjectFormValues = {
    name: string;
    slug: string;
    description: string;
    status: ProjectStatus;
    startDate: string;
    endDate: string;
    memberIds: number[];
};

type ProjectFormProps = {
    initialValues?: ProjectFormValues;
    users: CompanyUser[];
    isUsersLoading: boolean;
    usersErrorMessage: string;
    isSubmitting: boolean;
    submitErrorMessage: string;
    submitLabel: string;
    onSubmit: (
        request: ProjectWriteRequest,
    ) => Promise<void>;
    onCancel: () => void;
};

const defaultValues: ProjectFormValues = {
    name: '',
    slug: '',
    description: '',
    status: 'planned',
    startDate: '',
    endDate: '',
    memberIds: [],
};

function normalizeNullableValue(
    value: string,
): string | null {
    const normalizedValue = value.trim();

    return normalizedValue || null;
}

export function ProjectForm({
    initialValues = defaultValues,
    users,
    isUsersLoading,
    usersErrorMessage,
    isSubmitting,
    submitErrorMessage,
    submitLabel,
    onSubmit,
    onCancel,
}: ProjectFormProps) {
    const [values, setValues] =
        useState<ProjectFormValues>(
            initialValues,
        );

    const [validationMessage, setValidationMessage] =
        useState('');

    useEffect(() => {
        setValues(initialValues);
    }, [initialValues]);

    const selectedMemberIds = useMemo(
        () => new Set(values.memberIds),
        [values.memberIds],
    );

    function handleMemberChange(
        userId: number,
        checked: boolean,
    ): void {
        setValues((currentValues) => {
            const nextMemberIds = checked
                ? [
                      ...currentValues.memberIds,
                      userId,
                  ]
                : currentValues.memberIds.filter(
                      (memberId) =>
                          memberId !== userId,
                  );

            return {
                ...currentValues,
                memberIds: nextMemberIds,
            };
        });
    }

    async function handleSubmit(
        event: FormEvent<HTMLFormElement>,
    ): Promise<void> {
        event.preventDefault();
        setValidationMessage('');

        const name = values.name.trim();
        const slug = values.slug.trim();

        if (name.length < 2) {
            setValidationMessage(
                'プロジェクト名は2文字以上で入力してください。',
            );
            return;
        }

        if (slug.length < 2) {
            setValidationMessage(
                'スラッグは2文字以上で入力してください。',
            );
            return;
        }

        if (values.memberIds.length === 0) {
            setValidationMessage(
                'メンバーを1人以上選択してください。',
            );
            return;
        }

        if (
            values.startDate &&
            values.endDate &&
            values.endDate < values.startDate
        ) {
            setValidationMessage(
                '終了日は開始日以降の日付を指定してください。',
            );
            return;
        }

        await onSubmit({
            name,
            slug,
            description: normalizeNullableValue(
                values.description,
            ),
            status: values.status,
            start_date: normalizeNullableValue(
                values.startDate,
            ),
            end_date: normalizeNullableValue(
                values.endDate,
            ),
            member_ids: values.memberIds,
        });
    }

    return (
        <form
            className="project-form"
            onSubmit={(event) => {
                void handleSubmit(event);
            }}
        >
            {(validationMessage ||
                submitErrorMessage) && (
                <div
                    className="project-form__message project-form__message--error"
                    role="alert"
                >
                    {validationMessage ||
                        submitErrorMessage}
                </div>
            )}

            <div className="project-form__body">
                <div className="project-form__fields">
                    <label className="project-form__field">
                        <span className="project-form__label">
                            プロジェクト名
                            <span aria-hidden="true">
                                *
                            </span>
                        </span>

                        <input
                            className="project-form__input"
                            type="text"
                            required
                            minLength={2}
                            maxLength={150}
                            value={values.name}
                            onChange={(event) => {
                                setValues(
                                    (
                                        currentValues,
                                    ) => ({
                                        ...currentValues,
                                        name: event.target
                                            .value,
                                    }),
                                );
                            }}
                        />
                    </label>

                    <label className="project-form__field">
                        <span className="project-form__label">
                            スラッグ
                            <span aria-hidden="true">
                                *
                            </span>
                        </span>

                        <input
                            className="project-form__input"
                            type="text"
                            required
                            minLength={2}
                            maxLength={100}
                            value={values.slug}
                            onChange={(event) => {
                                setValues(
                                    (
                                        currentValues,
                                    ) => ({
                                        ...currentValues,
                                        slug: event.target
                                            .value,
                                    }),
                                );
                            }}
                        />

                        <span className="project-form__help">
                            半角英数字とハイフンの使用を推奨します。
                        </span>
                    </label>

                    <label className="project-form__field">
                        <span className="project-form__label">
                            説明
                        </span>

                        <textarea
                            className="project-form__textarea"
                            maxLength={5000}
                            rows={5}
                            value={values.description}
                            onChange={(event) => {
                                setValues(
                                    (
                                        currentValues,
                                    ) => ({
                                        ...currentValues,
                                        description:
                                            event.target
                                                .value,
                                    }),
                                );
                            }}
                        />
                    </label>

                    <label className="project-form__field">
                        <span className="project-form__label">
                            ステータス
                            <span aria-hidden="true">
                                *
                            </span>
                        </span>

                        <select
                            className="project-form__select"
                            value={values.status}
                            onChange={(event) => {
                                setValues(
                                    (
                                        currentValues,
                                    ) => ({
                                        ...currentValues,
                                        status: event
                                            .target
                                            .value as ProjectStatus,
                                    }),
                                );
                            }}
                        >
                            <option value="planned">
                                計画中
                            </option>
                            <option value="active">
                                進行中
                            </option>
                            <option value="completed">
                                完了
                            </option>
                            <option value="archived">
                                アーカイブ
                            </option>
                        </select>
                    </label>

                    <div className="project-form__date-fields">
                        <label className="project-form__field">
                            <span className="project-form__label">
                                開始日
                            </span>

                            <input
                                className="project-form__input"
                                type="date"
                                value={
                                    values.startDate
                                }
                                onChange={(event) => {
                                    setValues(
                                        (
                                            currentValues,
                                        ) => ({
                                            ...currentValues,
                                            startDate:
                                                event
                                                    .target
                                                    .value,
                                        }),
                                    );
                                }}
                            />
                        </label>

                        <label className="project-form__field">
                            <span className="project-form__label">
                                終了日
                            </span>

                            <input
                                className="project-form__input"
                                type="date"
                                min={
                                    values.startDate ||
                                    undefined
                                }
                                value={values.endDate}
                                onChange={(event) => {
                                    setValues(
                                        (
                                            currentValues,
                                        ) => ({
                                            ...currentValues,
                                            endDate:
                                                event
                                                    .target
                                                    .value,
                                        }),
                                    );
                                }}
                            />
                        </label>
                    </div>
                </div>

                <section className="project-form__members">
                    <h2 className="project-form__members-title">
                        メンバー
                        <span aria-hidden="true">
                            *
                        </span>
                    </h2>

                    {isUsersLoading ? (
                        <p
                            className="project-form__members-state"
                            role="status"
                        >
                            メンバーを読み込んでいます。
                        </p>
                    ) : usersErrorMessage ? (
                        <p
                            className="project-form__members-state project-form__members-state--error"
                            role="alert"
                        >
                            {usersErrorMessage}
                        </p>
                    ) : users.length === 0 ? (
                        <p className="project-form__members-state">
                            選択可能なメンバーがいません。
                        </p>
                    ) : (
                        <div className="project-form__member-list">
                            {users.map((user) => (
                                <label
                                    className="project-form__member"
                                    key={
                                        user.user_id
                                    }
                                >
                                    <input
                                        type="checkbox"
                                        checked={selectedMemberIds.has(
                                            user.user_id,
                                        )}
                                        onChange={(
                                            event,
                                        ) => {
                                            handleMemberChange(
                                                user.user_id,
                                                event
                                                    .target
                                                    .checked,
                                            );
                                        }}
                                    />

                                    <span className="project-form__member-information">
                                        <span className="project-form__member-name">
                                            {
                                                user.display_name
                                            }
                                        </span>

                                        <span className="project-form__member-email">
                                            {user.email}
                                        </span>
                                    </span>

                                    <span className="project-form__member-role">
                                        {user.role}
                                    </span>
                                </label>
                            ))}
                        </div>
                    )}
                </section>
            </div>

            <footer className="project-form__actions">
                <button
                    className="project-form__button project-form__button--secondary"
                    type="button"
                    disabled={isSubmitting}
                    onClick={onCancel}
                >
                    キャンセル
                </button>

                <button
                    className="project-form__button project-form__button--primary"
                    type="submit"
                    disabled={
                        isSubmitting ||
                        isUsersLoading ||
                        users.length === 0
                    }
                >
                    {isSubmitting
                        ? '処理中です。'
                        : submitLabel}
                </button>
            </footer>
        </form>
    );
}