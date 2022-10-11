def test_super_strnage_func():
    assert super_strnage_func(5) == 10


def test_validate_user_agent_chrome_good():
    assert basic_user_agent_validator("chrome 43.3.52.45.3") is True


def test_validate_user_agent_bad():
    assert basic_user_agent_validator("ch_rome 431.1.1.1222.") is False
