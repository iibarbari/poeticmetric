import { useRouter } from "next/router";
import React, { useCallback, useContext, useEffect, useRef, useState } from "react";
import { Alert, Button, Card, Container, Form, Spinner } from "react-bootstrap";
import { AuthAndApiContext, ToastsContext } from "../../contexts";
import { api, setUserAccessToken } from "../../helpers";
import { useForm } from "../../hooks";

type Form = {
  createDemoSite: boolean;
  organizationName: string;
  userEmail: string;
  userName: string;
  userNewPassword: string;
  userNewPassword2: string;
};

type State = {
  isDisabled: boolean;
  isReady: boolean;
};

export function Bootstrap() {
  const router = useRouter();
  const { mutate } = useContext(AuthAndApiContext);
  const { addToast } = useContext(ToastsContext);
  const isChecked = useRef<boolean>(false);
  const [state, setState] = useState<State>({ isDisabled: false, isReady: false });
  const [values, , updateValue, errors, setErrors] = useForm<Form>({
    createDemoSite: false,
    organizationName: "",
    userEmail: "",
    userName: "",
    userNewPassword: "",
    userNewPassword2: "",
  });

  const check = useCallback(async () => {
    const response = await api.get("/bootstrap-status");
    const responseJson = await response.json();

    if (!response.ok) {
      addToast({ body: responseJson.detail || "An error has occurred.", variant: "danger" });
      await router.replace("/");

      return;
    }

    if (!responseJson.isReady) {
      addToast({ body: "Bootstrap seems to be done already.", variant: "danger" });
      await router.replace("/");

      return;
    }

    setState((s) => ({ ...s, isReady: true }));
  }, [addToast, router]);

  const handleDemoSiteChange = useCallback<React.ChangeEventHandler<HTMLInputElement>>((event) => {
    updateValue("createDemoSite", event.target.checked);
  }, [updateValue]);

  const handleSubmit = useCallback<React.FormEventHandler<HTMLFormElement>>(async (event) => {
    event.preventDefault();

    setState((s) => ({ ...s, isDisabled: true }));

    const response = await api.post("/bootstrap", values);
    const responseJson = await response.json();

    if (!response.ok) {
      setErrors(responseJson);
      setState((s) => ({ ...s, isDisabled: false }));

      return;
    }

    setUserAccessToken(responseJson.userAccessToken.token);

    await mutate();

    addToast({ body: "Installation is complete!", variant: "success" });

    await router.replace("/sites");
  }, [addToast, mutate, router, setErrors, values]);

  useEffect(() => {
    if (!isChecked.current) {
      isChecked.current = true;

      check();
    }
  }, [check]);

  return state.isReady ? (
    <Container className="py-5">
      <div className="text-center">
        <h1>Welcome to PoeticMetric!</h1>

        <div className="mt-3">
          Initialize your PoeticMetric installation to continue.
        </div>
      </div>

      <Card className="mt-4 mx-auto mw-32rem">
        <Card.Body>
          {errors.detail !== undefined ? (
            <Alert variant="danger">{errors.detail}</Alert>
          ) : null}

          <Form onSubmit={handleSubmit}>
            <fieldset className="gap-3 vstack" disabled={state.isDisabled}>
              <Form.Group controlId="user-name">
                <Form.Label>Name</Form.Label>

                <Form.Control
                  isInvalid={errors.userName !== undefined}
                  maxLength={70}
                  minLength={1}
                  name="userName"
                  onChange={updateValue}
                  required
                />

                <Form.Control.Feedback type="invalid">{errors.userName}</Form.Control.Feedback>
              </Form.Group>

              <Form.Group controlId="email">
                <Form.Label>E-mail address</Form.Label>

                <Form.Control isInvalid={errors.userEmail !== undefined} name="userEmail" onChange={updateValue} required type="email" />

                <Form.Control.Feedback type="invalid">{errors.userEmail}</Form.Control.Feedback>
              </Form.Group>

              <Form.Group controlId="new-password">
                <Form.Label>New password</Form.Label>

                <Form.Control
                  isInvalid={errors.userNewPassword !== undefined}
                  name="userNewPassword"
                  onChange={updateValue}
                  required
                  type="password"
                />

                <Form.Control.Feedback type="invalid">{errors.userNewPassword}</Form.Control.Feedback>
              </Form.Group>

              <Form.Group controlId="new-password2">
                <Form.Label>New password (again)</Form.Label>

                <Form.Control
                  isInvalid={errors.userNewPassword2 !== undefined}
                  maxLength={72}
                  minLength={8}
                  name="userNewPassword2"
                  onChange={updateValue}
                  required
                  type="password"
                />

                <Form.Control.Feedback type="invalid">{errors.userNewPassword2}</Form.Control.Feedback>
              </Form.Group>

              <Form.Group controlId="organization-name">
                <Form.Label>Organization</Form.Label>

                <Form.Control
                  isInvalid={errors.organizationName !== undefined}
                  maxLength={70}
                  minLength={2}
                  name="organizationName"
                  onChange={updateValue}
                  required
                />

                <Form.Control.Feedback type="invalid">{errors.organizationName}</Form.Control.Feedback>
              </Form.Group>

              <Form.Group controlId="create-demo-site">
                <Form.Check
                  checked={values.createDemoSite}
                  label="Create demo site"
                  onChange={handleDemoSiteChange}
                />
              </Form.Group>

              <Button type="submit">Continue</Button>
            </fieldset>
          </Form>
        </Card.Body>
      </Card>
    </Container>
  ) : (
    <Spinner className="m-auto" variant="primary" />
  );
}
