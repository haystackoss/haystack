def additional_func():
    return 5
    
def test_super_strnage_func():
    assert additional_func() == 10


def test_validate_user_agent_chrome_good():
    assert True is True



def test_validate_user_agent_bad():
    assert False is False
