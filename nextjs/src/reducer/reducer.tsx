type State = {
  Image: string | null;
  Caption: string | null;
  Error: string | null;
};

export const initialState: State = {
  Image: null,
  Caption: null,
  Error: null,
};

type SetDisplayAction = {
  type: "SET_DISPLAY";
  payload: {
    Image: string;
    Caption: string;
  };
};

type ResetDisplayAction = {
  type: "RESET_DISPLAY";
};

type SetErrorAction = {
  type: "SET_ERROR";
  payload: {
    Error: string | null
  }
}

export type Action = SetDisplayAction | ResetDisplayAction | SetErrorAction;

export function PageReducer(
  state: State = initialState,
  action: Action,
): State {
  switch (action.type) {
    case "SET_DISPLAY":
      return {
        ...state,
        Image: action.payload.Image,
        Caption: action.payload.Caption,
      };
    case "RESET_DISPLAY":
      return initialState;
    case "SET_ERROR":
      return {
        ...initialState,
        Error: action.payload.Error
      }
    default:
      return state;
  }
}
